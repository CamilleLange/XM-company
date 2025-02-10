package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Aloe-Corporation/logs"
	"github.com/CamilleLange/XM-company/internal/config"
	"github.com/CamilleLange/XM-company/internal/features/company"
	"github.com/CamilleLange/XM-company/internal/interfaces/datasources"
	"github.com/CamilleLange/XM-company/internal/interfaces/events"
	ginhttp "github.com/CamilleLange/XM-company/internal/interfaces/http"
	"go.uber.org/zap"
)

var (
	log = logs.Get()
)

type RunCallback func()
type CloseCallback func() error

// Launch use the provided config to set up the API for running.
// This set up interfaces (datasources, http and event), shutdown operations, and finnaly launch the HTTP server.
func Launch(config config.Config) (RunCallback, CloseCallback, error) {
	// create datasources instances.
	mongo, err := datasources.NewMongoDB(config.Datasources.Mongo)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open connection to mongo : %w", err)
	}

	// create event handler instances.
	companyEvent, err := events.NewCompanyEventHandler(config.Event)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open connection to kafa : %w", err)
	}

	// create features.
	companyFeature, err := company.NewCompanyFeatures("mongo", mongo, companyEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create company feature : %w", err)
	}

	// create HTTP router.
	router := ginhttp.NewRouter(config.Router)

	// register each features endpoints.
	companyHandlers := ginhttp.NewCompanyHandler(companyFeature)
	companyHandlers.RegisterRoutes(router)

	// set up server.
	addrGin := config.Router.Addr + ":" + strconv.Itoa(config.Router.Port)
	srv := &http.Server{
		ReadHeaderTimeout: time.Millisecond,
		Addr:              addrGin,
		Handler:           router,
	}

	// prepare closing operations.
	close := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Router.ShutdownTimeout)*time.Second)
		defer cancel()

		if err := mongo.Disconnect(context.Background()); err != nil {
			return fmt.Errorf("can't disconnect mongo : %w", err)
		}

		if err := companyEvent.Close(); err != nil {
			return fmt.Errorf("can't close event handler : %w", err)
		}

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("can't shutdown server: %w", err)
		}

		return nil
	}

	// run the API.
	run := func() {
		log.Info("REST API listening on : "+addrGin,
			zap.String("package", "main"))

		log.Error(router.Run(addrGin).Error(),
			zap.String("package", "main"))
	}

	return run, close, nil
}
