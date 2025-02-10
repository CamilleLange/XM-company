package http

import (
	"time"

	"github.com/Aloe-Corporation/cors"
	"github.com/Aloe-Corporation/logs"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	// log is the singelton of the Zap logger.
	log = logs.Get()

	// ValidateInstance is the instance of the HTTP body structs validator.
	ValidateInstance *validator.Validate
)

// Config struct of the http package.
type Config struct {
	GinMode         string `mapstructure:"gin_mode"`
	Addr            string `mapstructure:"addr"`
	Port            int    `mapstructure:"port"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

// NewRouter create a *gin.Engine ready to add handlers.
func NewRouter(config Config) *gin.Engine {
	ValidateInstance = validator.New()

	router := gin.New()
	gin.SetMode(config.GinMode)

	router.Use(ginzap.RecoveryWithZap(log, true))
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(cors.Middleware(nil))

	return router
}
