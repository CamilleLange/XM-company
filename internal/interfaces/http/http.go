package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Aloe-Corporation/cors"
	"github.com/Aloe-Corporation/logs"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var (
	// log is the singelton of the Zap logger.
	log = logs.Get()

	// ValidateInstance is the instance of the HTTP body structs validator.
	ValidateInstance *validator.Validate

	// jwtConfig store the configuration for the jwt middleware.
	jwtConfig ConfigAuth
)

// Config struct of the http package.
type Config struct {
	GinMode         string     `mapstructure:"gin_mode"`
	Addr            string     `mapstructure:"addr"`
	Port            int        `mapstructure:"port"`
	ShutdownTimeout int        `mapstructure:"shutdown_timeout"`
	Auth            ConfigAuth `mapstructure:"auth"`
}

type ConfigAuth struct {
	JWTKey   string `mapstructure:"jwt_key"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// NewRouter create a *gin.Engine ready to add handlers.
func NewRouter(config Config) *gin.Engine {
	jwtConfig = config.Auth

	ValidateInstance = validator.New()

	router := gin.New()
	gin.SetMode(config.GinMode)

	router.Use(ginzap.RecoveryWithZap(log, true))
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(cors.Middleware(nil))

	router.POST("/login", Login)

	return router
}

// Login authenticate the user with a JWT.
func Login(c *gin.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		log.Info("fail to parse credentials", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if credentials.Username != jwtConfig.Username || credentials.Password != jwtConfig.Password {
		log.Info("invalid credentials",
			zap.String("username", credentials.Username),
			zap.String("password", credentials.Password),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := GenerateJWT(credentials.Username)
	if err != nil {
		log.Warn("fail to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// GenerateJWT for the provided username.
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := token.SignedString([]byte(jwtConfig.JWTKey))
	if err != nil {
		return "", fmt.Errorf("cant sign token : %w", err)
	}
	return signedString, nil
}

// BasicJWTMiddleware handle JWT verification on protected endpoints.
func BasicJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Info("no authorization header provided. Aborting request.")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		tokenString := authHeader[len("Bearer "):]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.JWTKey), nil
		})

		if err != nil || !token.Valid {
			log.Warn("invalid or expired token. Aborting request.")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
