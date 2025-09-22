package main

import (
	"context"
	"errors"
	"main/internal/config"
	"main/internal/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

func main() {
	conf := config.Load(os.Getenv("CONFIG_PATH"))

	localLogger := logger.New("prod")

	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/swagger.json", func(c *gin.Context) {
		b, err := os.ReadFile(conf.OpenApiPath)
		if err != nil {
			localLogger.Errorw("failed to read swagger.json", "err", err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.Data(http.StatusOK, "application/json", b)
	})

	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	)))

	addr := "0.0.0.0:" + conf.Handler.Swagger.Port
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	localLogger.Infow("Listening", "port", conf.Handler.Swagger.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			localLogger.Errorw("HTTP server failed", "err", err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	localLogger.Infow("waiting for shutdown")

	<-ch

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		localLogger.Errorw("server shutdown error", zap.Error(err))
	} else {
		localLogger.Infow("shutdown complete")
	}
}
