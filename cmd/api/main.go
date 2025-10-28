package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"backend-challenge/internal/application"
	"backend-challenge/internal/config"
	jwtinfra "backend-challenge/internal/infrastructure/jwt"
	mongorepo "backend-challenge/internal/infrastructure/mongo"
	grpcsvc "backend-challenge/internal/transport/grpcsvc"
	transport "backend-challenge/internal/transport/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	httpHandler := transport.NewHandler(userService, jwtManager)
	httpRouter := transport.NewRouter(httpHandler, jwtManager)
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: httpRouter,
	}

	grpcListener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen grpc: %v", err)
	}
	defer grpcListener.Close()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcsvc.AuthUnaryInterceptor(jwtManager)),
	)
	grpcService := grpcsvc.NewUserServer(userService, jwtManager)
	grpcService.Register(grpcServer)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	group, groupCtx := errgroup.WithContext(ctx)
	shutdownOnce := &sync.Once{}

	shutdown := func(reason string) {
		shutdownOnce.Do(func() {
			log.Printf("shutdown initiated: %s", reason)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			if err := httpServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("http shutdown error: %v", err)
			}

			stopped := make(chan struct{})
			go func() {
				grpcServer.GracefulStop()
				close(stopped)
			}()

			select {
			case <-stopped:
			case <-time.After(5 * time.Second):
				log.Println("forcing grpc stop")
				grpcServer.Stop()
			}
		})
	}

	go func() {
		<-ctx.Done()
		shutdown("signal")
	}()

	group.Go(func() error {
		log.Printf("http server listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			shutdown("http server error")
			return err
		}
		return nil
	})

	group.Go(func() error {
		log.Printf("grpc server listening on %s", grpcListener.Addr().String())
		if err := grpcServer.Serve(grpcListener); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			if status.Code(err) == codes.Canceled {
				return nil
			}
			shutdown("grpc server error")
			return err
		}
		return nil
	})

	group.Go(func() error {
		runUserCountWorker(groupCtx, userService, cfg.BackgroundTick)
		return nil
	})

	if err := group.Wait(); err != nil {
		log.Printf("server stopped with error: %v", err)
	} else {
		log.Println("server stopped gracefully")
	}
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
