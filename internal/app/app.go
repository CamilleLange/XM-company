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
	mongo, err := datasources.NewMongoDB(config.Datasources.Mongo)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open connection to mongo : %w", err)
	}

	companyFeature, err := company.NewCompanyFeatures("mongo", mongo)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create company feature : %w", err)
	}

	router := ginhttp.NewRouter(config.Router)

	companyHandlers := ginhttp.NewCompanyHandler(companyFeature)
	companyHandlers.RegisterRoutes(router)

	addrGin := config.Router.Addr + ":" + strconv.Itoa(config.Router.Port)
	srv := &http.Server{
		ReadHeaderTimeout: time.Millisecond,
		Addr:              addrGin,
		Handler:           router,
	}

	close := func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Router.ShutdownTimeout)*time.Second)
		defer cancel()

		if err := mongo.Disconnect(context.Background()); err != nil {
			return fmt.Errorf("can't disconnect mongo : %w", err)
		}

		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("can't shutdown server: %w", err)
		}

		return nil
	}

	run := func() {
		log.Info("REST API listening on : "+addrGin,
			zap.String("package", "main"))

		log.Error(router.Run(addrGin).Error(),
			zap.String("package", "main"))
	}

	return run, close, nil
}
