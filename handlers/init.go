package handlers

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"

	"togo/config"
	"togo/logger"
	"togo/mongodb"
)

var (
	log = logger.GetLogger("handlers")
	e   = echo.New()

	mongoClient mongodb.MongoStore
)

// CustomValidator ...
type CustomValidator struct {
	validator *validator.Validate
}

// Validate ...
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func initClient() (err error) {
	mongoClient, err = mongodb.NewMongoDBClient()
	if err != nil {
		log.Errorf("Connect Mongo failed: %v\n", err)
		return
	}

	return
}

func stopClients() {
	log.Info("Stopping all clients")

	mongoClient.Close()

	log.Info("All clients have been closed")
}

// healthCheck handler
func healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func startServer() {
	err := initClient()
	if err != nil {
		stopClients()
		log.Panicf("fatal error initClient: %s", err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Validator = &CustomValidator{validator: validator.New()}

	// public apis
	e.GET("/health", healthCheck)
	e.POST("/user/register", registerUser)
	e.POST("/user/login", login)

	// private apis
	privateAPIGroup := e.Group("/private")
	privateAPIGroup.Use(isLoggedIn, private)
	privateAPIGroup.POST("/task/add", addTask)

	log.Info("Starting at port: " + config.Values.Port)
	if err := e.Start(":" + config.Values.Port); err != nil {
		log.Error(err)
	}
}

func waitForInterruptSignal() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Error(err)
	}

	time.Sleep(5 * time.Second)
	stopClients()
}

// Start APIs service
func Start() {
	go startServer()
	waitForInterruptSignal()
}
