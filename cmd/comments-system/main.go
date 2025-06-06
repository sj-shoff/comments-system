package main

import (
	"comments-system/internal/config"
	"comments-system/internal/graph"
	"comments-system/internal/graph/generated"
	"comments-system/internal/pubsub"
	"comments-system/internal/service"
	"comments-system/internal/storage"
	"comments-system/internal/storage/inmemory"
	"comments-system/internal/storage/postgres"
	"comments-system/pkg/logger/sl"
	"comments-system/pkg/logger/slogpretty"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"golang.org/x/sync/errgroup"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("Starting server", "env", cfg.Env, "storage", cfg.Storage)

	var storage storage.Storage
	var err error

	switch cfg.Storage {
	case "postgres":
		storage, err = postgres.NewPostgresDB(cfg.Database)
		if err != nil {
			log.Error("Failed to init postgres", sl.Err(err))
			os.Exit(1)
		}
		log.Info("Using PostgreSQL storage")
	case "inmemory":
		storage = inmemory.NewInMemory()
		log.Info("Using in-memory storage")
	}

	postService := service.NewPostService(storage, log)
	commentService := service.NewCommentService(storage, log)
	services := &service.Service{
		PostService:    postService,
		CommentService: commentService,
	}

	ps := pubsub.NewPubSub()

	srv := handler.New(generated.NewExecutableSchema(generated.Config{
		Resolvers: graph.NewResolver(services, ps, log),
	}))

	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	srv.Use(extension.Introspection{})

	router := http.NewServeMux()
	router.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	router.Handle("/query", graph.ContentTypeMiddleware(srv))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Info("Server listening on", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		<-gCtx.Done()
		log.Info("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}
		log.Info("Server stopped")
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Error("Server terminated with error", sl.Err(err))
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
