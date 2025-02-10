package app

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Aloe-Corporation/logs"
	"github.com/CamilleLange/XM-compagny/internal/config"
	"github.com/CamilleLange/XM-compagny/internal/features/compagny"
	"github.com/CamilleLange/XM-compagny/internal/interfaces/datasources"
	ginhttp "github.com/CamilleLange/XM-compagny/internal/interfaces/http"
	"go.uber.org/zap"
)

var (
	log = logs.Get()
)

type RunCallback func()
type CloseCallback func() error

func Launch(config config.Config) (RunCallback, CloseCallback, error) {
	mongo, err := datasources.NewMongoDB(config.Datasources.Mongo)
	if err != nil {
		return nil, nil, fmt.Errorf("can't open connection to mongo : %w", err)
	}

	compagnyFeature, err := compagny.NewCompagnyFeatures("mongo", mongo)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create compagny feature : %w", err)
	}

	router := ginhttp.NewRouter(config.Router)

	compagnyHandlers := ginhttp.NewCompagnyHandler(compagnyFeature)
	compagnyHandlers.RegisterRoutes(router)

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
