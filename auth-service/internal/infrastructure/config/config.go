package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type Config struct {
	App      App      `yaml:"app" validate:"required"`
	Server   Server   `yaml:"server" validate:"required"`
	Database Database `yaml:"database" validate:"required"`
	Cache    Cache    `yaml:"cache" validate:"required"`
	Session  Session  `yaml:"session" validate:"required"`
}

type App struct {
	Name    string `yaml:"name" validate:"required"`
	Version string `yaml:"version" validate:"required,semver"`
	Env     string `yaml:"env" validate:"required,oneof=prod dev local"`
}

type Server struct {
	GRPC struct {
		Addr string `yaml:"addr" validate:"required,hostname_port"`
		MTLS struct {
			Enable         bool   `yaml:"enable"`
			ServerCertPath string `yaml:"server_cert_path"`
			ServerKeyPath  string `yaml:"server_key_path"`
			RootCAPath     string `yaml:"root_ca_path"`
		} `yaml:"mtls"`
	} `yaml:"grpc" validate:"required"`
	Timeout time.Duration `yaml:"timeout" validate:"required,min=100ms"`
}

type Database struct {
	Postgres struct {
		Host    string `yaml:"host" validate:"required,hostname"`
		Port    int    `yaml:"port" validate:"required,gte=1,lte=65535"`
		DBName  string `yaml:"dbname" validate:"required"`
		SSLMode string `yaml:"sslmode" validate:"required,oneof=disable enable"`
		Auth    struct {
			User     string `yaml:"user" validate:"required"`
			Password string `yaml:"password" validate:"required"`
		} `yaml:"auth"`
		Conn struct {
			MaxOpens    int           `yaml:"max_opens"`
			MaxIdles    int           `yaml:"max_idles"`
			MaxLifetime time.Duration `yaml:"max_lifetime"`
			MaxIdleTime time.Duration `yaml:"max_idle_time"`
		} `yaml:"conn"`
	} `yaml:"postgres"`
}

type Cache struct {
	Redis struct {
		Host string `yaml:"host" validate:"required,hostname"`
		Port int    `yaml:"port" validate:"required,gte=1,lte=65535"`
		DB   int    `yaml:"db"`
		Auth struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"auth"`
	} `yaml:"redis"`
}

type Session struct {
	RefreshTTL     time.Duration `yaml:"refresh_ttl" validate:"required,min=100ms"`
	AccessTTL      time.Duration `yaml:"access_ttl" validate:"required,min=100ms"`
	PublicKeyPath  string        `yaml:"public_key_path" validate:"required,file"`
	PrivateKeyPath string        `yaml:"private_key_path" validate:"required,file"`
}

func New(path string) (*Config, error) {
	bytes, err := parseFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg, err := parseBytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return cfg, nil
}

func validateConfig(cfg *Config) error {
	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		return err
	}

	if cfg.Server.GRPC.MTLS.Enable {
		if _, err := os.Stat(cfg.Server.GRPC.MTLS.ServerCertPath); os.IsNotExist(err) {
			return errors.New("server cert file doesn't exist")
		}

		if _, err := os.Stat(cfg.Server.GRPC.MTLS.ServerKeyPath); os.IsNotExist(err) {
			return errors.New("server key file doesn't exist")
		}

		if _, err := os.Stat(cfg.Server.GRPC.MTLS.RootCAPath); os.IsNotExist(err) {
			return errors.New("root ca doesn't exist")
		}
	}

	return nil
}

func parseBytes(bytes []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseFile(path string) ([]byte, error) {
	path = filepath.Clean(path)
	if _, err := filepath.Abs(path); err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := os.ExpandEnv(string(bytes))

	return []byte(content), nil
}
