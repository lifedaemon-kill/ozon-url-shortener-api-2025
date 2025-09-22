package main

import (
	"context"
	"errors"
	"main/internal/config"
	"main/internal/logger"
	desc "main/pkg/protogen/url-shortener"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.Load(os.Getenv("CONFIG_PATH"))

	localLogger := logger.New(cfg.Env)

	// gRPC-Gateway mux
	gwMux := runtime.NewServeMux()
	err := desc.RegisterUrlShortenerHandlerFromEndpoint(
		ctx, gwMux, "grpc-api:"+cfg.Handler.Grpc.Port,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		localLogger.Errorw("failed to register gRPC-Gateway handler", "err", err)
		cancel()
	}

	// Gin engine
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // или конкретно: []string{"http://localhost:8081"}
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	gin.SetMode(gin.ReleaseMode)

	// Gin -> wrap gRPC-Gateway ServeMux
	r.Any("/*any", gin.WrapH(gwMux))

	srv := &http.Server{
		Addr:    "0.0.0.0:" + cfg.Handler.HttpGW.Port,
		Handler: r,
	}

	localLogger.Infow("Starting HTTP server", "port ", cfg.Handler.HttpGW.Port)
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			localLogger.Errorw("server error", "err", err)
			cancel()
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	localLogger.Infow("waiting for shutdown")

	<-ch

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		localLogger.Errorw("server shutdown error", zap.Error(err))
	} else {
		localLogger.Infow("shutdown complete")
	}
}
