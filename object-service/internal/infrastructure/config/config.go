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
	Storage  Storage  `yaml:"storage" validate:"required"`
}

type App struct {
	Name    string `yaml:"name" validate:"required"`
	Version string `yaml:"version" validate:"required,semver"`
	Env     string `yaml:"env" validate:"required,oneof=dev prod local"`
}

type Server struct {
	HTTP struct {
		Addr string `yaml:"addr" validate:"required,hostname_port"`
		TLS  struct {
			Enable         bool   `yaml:"enable"`
			ServerCertPath string `yaml:"server_cert_path"`
			ServerKeyPath  string `yaml:"server_key_path"`
		} `yaml:"tls"`
	} `yaml:"http" validate:"required"`
	ReadTimeout  time.Duration `yaml:"read_timeout" validate:"required,min=100ms"`
	WriteTimeout time.Duration `yaml:"write_timeout" validate:"required,min=100ms"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" validate:"required,min=100ms"`
}

type Database struct {
	Postgres struct {
		Host   string `yaml:"host" validate:"required,hostname"`
		Port   int    `yaml:"port" validate:"required,gte=1,lte=65535"`
		DBName string `yaml:"dbname" validate:"required"`
		Auth   struct {
			User     string `yaml:"user"`
			Password string `yaml:"password"`
		} `yaml:"auth"`
		Conn struct {
			MaxIdles    int           `yaml:"max_idles" validate:"required,gte=1"`
			MaxOpens    int           `yaml:"max_opens" validate:"required,gte=1"`
			MaxIdleTime time.Duration `yaml:"max_idle_time" validate:"required,min=100ms"`
			MaxLifetime time.Duration `yaml:"max_lifetime" validate:"required,min=100ms"`
		} `yaml:"conn"`
	} `yaml:"postgres"`
}

type Storage struct {
	BasePath string `yaml:"base_path" validate:"required,path"`
}

func New(path string) (*Config, error) {
	bytes, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg, err := parseBytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bytes: %w", err)
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

	if cfg.Server.HTTP.TLS.Enable {
		if _, err := os.Stat(cfg.Server.HTTP.TLS.ServerCertPath); os.IsNotExist(err) {
			return errors.New("server cert file doesn't exist")
		}
		if _, err := os.Stat(cfg.Server.HTTP.TLS.ServerKeyPath); os.IsNotExist(err) {
			return errors.New("server key file doesn't exist")
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

func readFile(path string) ([]byte, error) {
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
