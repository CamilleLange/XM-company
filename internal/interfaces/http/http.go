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
	log              = logs.Get()
	ValidateInstance *validator.Validate
)

type Config struct {
	GinMode         string `mapstructure:"gin_mode"`
	Addr            string `mapstructure:"addr"`
	Port            int    `mapstructure:"port"`
	ShutdownTimeout int    `mapstructure:"shutdown_timeout"`
}

func NewRouter(config Config) *gin.Engine {
	ValidateInstance = validator.New()

	router := gin.New()
	gin.SetMode(config.GinMode)

	router.Use(ginzap.RecoveryWithZap(log, true))
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(cors.Middleware(nil))

	return router
}
