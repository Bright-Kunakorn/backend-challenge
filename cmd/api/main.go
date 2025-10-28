package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend-challenge/internal/application"
	"backend-challenge/internal/config"
	jwtinfra "backend-challenge/internal/infrastructure/jwt"
	mongorepo "backend-challenge/internal/infrastructure/mongo"
	transport "backend-challenge/internal/transport/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("connect to mongo: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("ping mongo: %v", err)
	}

	db := client.Database(cfg.MongoDatabase)
	userRepo, err := mongorepo.NewUserRepository(db)
	if err != nil {
		log.Fatalf("init user repository: %v", err)
	}

	userService := application.NewUserService(userRepo)
	jwtManager := jwtinfra.NewManager(cfg.JWTSecret, cfg.JWTExpiry, cfg.JWTIssuer)
	handler := transport.NewHandler(userService, jwtManager)
	router := transport.NewRouter(handler, jwtManager)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			serverErrors <- err
		}
	}()

	workerCtx, workerCancel := context.WithCancel(context.Background())
	defer workerCancel()
	go runUserCountWorker(workerCtx, userService, cfg.BackgroundTick)

	select {
	case err := <-serverErrors:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	case <-ctx.Done():
		log.Println("shutdown signal received")
	}

	workerCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}

	log.Println("server stopped")
}

func runUserCountWorker(ctx context.Context, service *application.UserService, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			countCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			count, err := service.Count(countCtx)
			cancel()
			if err != nil {
				log.Printf("user count worker error: %v", err)
				continue
			}
			log.Printf("user count: %d", count)
		}
	}
}
