package config

import (
	"errors"
	"fmt"
	"log"
	"main/internal/errs"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string        `yaml:"env"`
	Handler      HandlerConfig `yaml:"handler"`
	Storage      StorageConfig `yaml:"storage"`
	DB           DBConfig      `yaml:"db"`
	UrlGenerator UrlGenConfig  `yaml:"url-generator"`
	OpenApiPath  string        `yaml:"openapi-path"`
}
type UrlGenConfig struct {
	Alphabet string `yaml:"alphabet"`
	Length   int    `yaml:"length"`
	BaseHost string `yaml:"base-host"`
}
type DBConfig struct {
	Host          string `env:"DB_HOST"     env-default:"postgres"`
	Port          int    `env:"DB_PORT"     env-default:"5432"`
	User          string `env:"POSTGRES_USER"     env-default:"postgres"`
	Password      string `env:"POSTGRES_PASSWORD" env-required:"true"`
	Name          string `env:"POSTGRES_DB"     env-default:"ozon_db"`
	SSLMode       string `env:"DB_SSLMODE"  env-default:"disable"`
	MigrationPath string `yaml:"migration-path"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
	)
}

type StorageConfig struct {
	Type string `yaml:"type" required:"true" default:"postgres"`
}

func Load(path string) Config {
	conf := Config{}
	if err := cleanenv.ReadConfig(path, &conf); err != nil {
		log.Fatal(errors.Join(errs.ErrorConfigFileNotFound, errors.New("not found: "+path), err))
	}
	if err := cleanenv.ReadEnv(&conf); err != nil {
		log.Fatal(errors.Join(errs.ErrorConfigFileNotFound, errors.New("not found: "+path), err))
	}

	return conf
}

type HandlerConfig struct {
	Grpc    GrpcConfig    `yaml:"grpc"`
	HttpGW  HttpGWConfig  `yaml:"http-gw"`
	Swagger SwaggerConfig `yaml:"swagger"`
}

type GrpcConfig struct {
	Port string `yaml:"port" required:"true"`
}
type HttpGWConfig struct {
	Port string `yaml:"port" required:"true"`
}
type SwaggerConfig struct {
	Port string `yaml:"port" required:"true"`
}
