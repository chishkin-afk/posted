package postgres

import (
	"fmt"

	"github.com/chishkin-afk/posted/posts-service/internal/infrastructure/config"
	postpg "github.com/chishkin-afk/posted/posts-service/internal/infrastructure/persistence/postgres/post"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d sslmode=%s user=%s password=%s dbname=%s",
		cfg.Database.Postgres.Host,
		cfg.Database.Postgres.Port,
		cfg.Database.Postgres.SSLMode,
		cfg.Database.Postgres.Auth.User,
		cfg.Database.Postgres.Auth.Password,
		cfg.Database.Postgres.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, fmt.Errorf("failed to open connection with postgres: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to open sql db: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.Postgres.Conn.MaxOpens)
	sqlDB.SetMaxIdleConns(cfg.Database.Postgres.Conn.MaxIdles)
	sqlDB.SetConnMaxIdleTime(cfg.Database.Postgres.Conn.MaxIdleTime)
	sqlDB.SetConnMaxLifetime(cfg.Database.Postgres.Conn.MaxLifetime)

	return db, nil
}

func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql connection: %w", err)
	}

	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close sql connection: %w", err)
	}

	return nil
}

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&postpg.PostModel{}); err != nil {
		return fmt.Errorf("failed to migrate db: %w", err)
	}

	return nil
}
