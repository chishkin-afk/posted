package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	authpb "github.com/chishkin-afk/posted/http-gateway/api/auth/v1"
	postspb "github.com/chishkin-afk/posted/http-gateway/api/posts/v1"
	"github.com/chishkin-afk/posted/http-gateway/internal/application/services"
	"github.com/chishkin-afk/posted/http-gateway/internal/infrastructure/config"
	httpserver "github.com/chishkin-afk/posted/http-gateway/internal/infrastructure/http"
	"github.com/chishkin-afk/posted/http-gateway/internal/infrastructure/http/handlers"
	"github.com/chishkin-afk/posted/http-gateway/internal/infrastructure/mtls"
	"github.com/chishkin-afk/posted/http-gateway/pkg/log"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	log    *slog.Logger
	server *httpserver.Server
}

func (a *App) Start() error {
	a.log.Info("server is running")
	return a.server.Start()
}

func (a *App) Shutdown(ctx context.Context) error {
	a.log.Info("server shutdown")
	return a.server.Shutdown(ctx)
}

func New() (*App, func(), error) {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("failed to load .env",
			slog.String("error", err.Error()),
		)
	}

	cfg, err := config.New(os.Getenv("APP_CONFIG_PATH"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	log := log.New(cfg.App.Env)

	log.Info("config was loaded", slog.Any("server", cfg.Server))

	authClientConn, err := provideExternalService(&cfg.GRPC.AuthService)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide auth service: %w", err)
	}

	postsClientConn, err := provideExternalService(&cfg.GRPC.PostsService)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide posts service: %w", err)
	}

	server, err := provideServer(cfg, log, authClientConn, postsClientConn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide server: %w", err)
	}

	return &App{
			log:    log,
			server: server,
		}, func() {
			if err := authClientConn.Close(); err != nil {
				log.Error("failed to close connection with auth", slog.String("error", err.Error()))
			} else {
				log.Info("connection with auth was closed")
			}

			if err := postsClientConn.Close(); err != nil {
				log.Error("failed to close connection with posts", slog.String("error", err.Error()))
			} else {
				log.Info("connection with posts was closed")
			}
		}, nil
}

func provideServer(cfg *config.Config, log *slog.Logger, authCC, postsCC *grpc.ClientConn) (*httpserver.Server, error) {
	authService := services.NewAuthService(
		cfg,
		log,
		authpb.NewAuthServiceClient(authCC),
	)
	postsService := services.NewPostsService(
		cfg,
		log,
		postspb.NewPostsServiceClient(postsCC),
	)

	handlers, err := handlers.New(
		cfg.App.Env,
		authService,
		postsService,
	)
	if err != nil {
		return nil, err
	}

	return httpserver.New(cfg, handlers), nil
}

func provideExternalService(cfg *config.ExternalServer) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}
	if cfg.MTLS.Enable {
		tlsConfig, err := mtls.LoadMTLS(&cfg.MTLS)
		if err != nil {
			return nil, fmt.Errorf("failed to provide auth service: %w", err)
		}

		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	clientConn, err := grpc.NewClient(cfg.Addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to open client conn: %w", err)
	}

	return clientConn, nil
}
