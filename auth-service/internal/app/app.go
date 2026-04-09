package app

import (
	"fmt"
	"log/slog"
	"os"

	authpb "github.com/chishkin-afk/posted/auth-service/api/auth/v1"
	"github.com/chishkin-afk/posted/auth-service/internal/application/services"
	"github.com/chishkin-afk/posted/auth-service/internal/domain/session"
	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/config"
	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/grpc/server"
	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/grpc/server/handlers"
	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/grpc/server/interceptors"
	jwt "github.com/chishkin-afk/posted/auth-service/internal/infrastructure/jwt/session"
	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/mtls"
	"github.com/chishkin-afk/posted/auth-service/internal/infrastructure/persistence/postgres"
	userpg "github.com/chishkin-afk/posted/auth-service/internal/infrastructure/persistence/postgres/user"
	"github.com/chishkin-afk/posted/auth-service/pkg/log"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"
)

type App struct {
	server *server.Server
	log    *slog.Logger
}

func (a *App) Start() error {
	a.log.Info("server is running")
	return a.server.Start()
}

func (a *App) GracefulStop() {
	a.log.Info("server shutdown")
	a.server.GracefulStop()
}

func New() (*App, func(), error) {
	if err := godotenv.Load(".env"); err != nil {
		slog.Warn("failed to load .env", slog.String("error", err.Error()))
	}

	cfg, err := config.New(os.Getenv("APP_CONFIG_PATH"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	log := log.New(cfg.App.Env)

	log.Info("config was loaded", slog.Any("server", cfg.Server))

	db, err := providePersistence(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide persistence: %w", err)
	}

	jm, err := jwt.New(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load jwt manager: %w", err)
	}

	server, err := provideServer(cfg, log, db, jm)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide server: %w", err)
	}

	return &App{
			server: server,
			log:    log,
		}, func() {
			if err := postgres.Close(db); err != nil {
				log.Warn("failed to close connection with postgres", slog.String("error", err.Error()))
			} else {
				log.Info("connection with postgres was closed")
			}
		}, nil
}

func provideServer(cfg *config.Config, log *slog.Logger, db *gorm.DB, jm session.JWTManager) (*server.Server, error) {
	service := services.New(
		cfg,
		log,
		userpg.New(db),
		jm,
	)

	handler := handlers.New(service)
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.NewAuthInterceptor(jm, map[string]bool{
				authpb.AuthService_DeleteUser_FullMethodName:  true,
				authpb.AuthService_GetUserSelf_FullMethodName: true,
				authpb.AuthService_UpdateUser_FullMethodName:  true,
			}),
		),
	}
	if cfg.Server.GRPC.MTLS.Enable {
		tlsConfig, err := mtls.LoadMTLSConfig(cfg)
		if err != nil {
			return nil, err
		}

		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}

	grpcServer := grpc.NewServer(opts...)
	authpb.RegisterAuthServiceServer(grpcServer, handler)

	return server.New(cfg, grpcServer), nil
}

func providePersistence(cfg *config.Config) (*gorm.DB, error) {
	db, err := postgres.Connect(cfg)
	if err != nil {
		return nil, err
	}

	if err := postgres.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
