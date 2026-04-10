package app

import (
	"fmt"
	"log/slog"
	"os"

	postspb "github.com/chishkin-afk/posted/posts-service/api/posts/v1"
	"github.com/chishkin-afk/posted/posts-service/internal/application/services"
	"github.com/chishkin-afk/posted/posts-service/internal/domain/session"
	redissdk "github.com/chishkin-afk/posted/posts-service/internal/infrastructure/cache/redis"
	postredis "github.com/chishkin-afk/posted/posts-service/internal/infrastructure/cache/redis/post"
	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/config"
	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/grpc/server"
	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/grpc/server/handlers"
	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/grpc/server/interceptors"
	jwt "github.com/chishkin-afk/posted/posts-service/internal/infrastructure/jwt/session"
	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/mtls"
	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/persistence/postgres"
	postpg "github.com/chishkin-afk/posted/posts-service/internal/infrastructure/persistence/postgres/post"
	"github.com/chishkin-afk/posted/posts-service/pkg/log"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"

	"buf.build/go/protovalidate"

	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
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

	db, err := provideDatabase(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide database: %w", err)
	}

	redisClient, err := redissdk.Connect(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	jm, err := jwt.New(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load jwt manager: %w", err)
	}

	server, err := provideServer(cfg, log, db, redisClient, jm)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to provide server: %w", err)
	}

	return &App{
			log:    log,
			server: server,
		}, func() {
			if err := postgres.Close(db); err != nil {
				log.Error("failed to close postgres connection", slog.String("error", err.Error()))
			} else {
				log.Info("connection with postgres was closed")
			}

			if err := redissdk.Close(redisClient); err != nil {
				log.Error("failed to close redis connection", slog.String("error", err.Error()))
			} else {
				log.Info("connection with redis was closed")
			}
		}, nil
}

func provideServer(
	cfg *config.Config,
	log *slog.Logger,
	db *gorm.DB,
	redisClient *redis.Client,
	jm session.JWTManager,
) (*server.Server, error) {
	service := services.New(
		cfg,
		log,
		postpg.New(db),
		postredis.New(cfg, redisClient),
		jm,
	)

	handler := handlers.New(service)

	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.NewAuthInterceptor(jm, map[string]bool{
				postspb.PostsService_Create_FullMethodName:       true,
				postspb.PostsService_Delete_FullMethodName:       true,
				postspb.PostsService_GetSelfPosts_FullMethodName: true,
				postspb.PostsService_Update_FullMethodName:       true,
			}),
			protovalidate_middleware.UnaryServerInterceptor(validator),
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
	postspb.RegisterPostsServiceServer(grpcServer, handler)

	return server.New(cfg, grpcServer), nil
}

func provideDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := postgres.Connect(cfg)
	if err != nil {
		return nil, err
	}

	if err := postgres.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
