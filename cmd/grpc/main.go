package main

import (
	"context"
	"log"
	"main/internal/config"
	"main/internal/generator"
	grpcHandler "main/internal/handler/grpc"
	"main/internal/logger"
	"main/internal/repository"
	"main/internal/repository/inmemory"
	"main/internal/repository/postgres"
	"main/internal/service"
	"main/pkg/cache"
	"main/pkg/db/migration"
	desc "main/pkg/protogen/url-shortener"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conf := config.Load(os.Getenv("CONFIG_PATH"))
	logSugar := logger.New(conf.Env)

	logSugar.Infow("HELLO", "CONFIG", conf)
	var db repository.UrlRepository

	switch conf.Storage.Type {
	case "postgres":
		logSugar.Infow("using postgres", "dsn", conf.DB.DSN())
		var err error
		db, err = postgres.New(conf.DB)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		if err = migration.Run("migrations/", conf.DB.DSN()); err != nil {
			log.Fatal(err)
		}

	case "inmemory":
		db = inmemory.New()
	default:
		log.Fatal("wrong storage type: " + conf.Storage.Type)
	}

	urlGenerator := generator.New(conf.UrlGenerator)
	lfuCache := cache.NewLFUCache[string, string](60*time.Minute, 20)
	defer lfuCache.Stop()
	urlService := service.NewUrlService(urlGenerator, db, lfuCache, logSugar)

	lis, err := net.Listen("tcp", "0.0.0.0:"+conf.Handler.Grpc.Port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	grpcH := grpcHandler.New(urlService, logSugar)
	desc.RegisterUrlShortenerServer(grpcServer, grpcH)

	logSugar.Infow("Listening gRPC", "port", conf.Handler.Grpc.Port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			logSugar.Errorw("failed to serve", "err", err)
			return
		}
	}()
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	logSugar.Infow("shutting down gracefully...")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		logSugar.Infow("gRPC server stopped")
	case <-time.After(5 * time.Second):
		logSugar.Infow("timeout, forcing stop")
		grpcServer.Stop()
	}
}
