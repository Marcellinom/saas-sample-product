package config

import (
	"log"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/samber/do"
)

type AppConfig struct {
	Name  string `env:"APP_NAME,default=dptsi-base-go"`
	Env   string `env:"APP_ENV,default=production"`
	Key   string `env:"APP_KEY"`
	Debug bool   `env:"APP_DEBUG,default=false"`
	URL   string `env:"APP_URL,default=http://localhost"`
}

type CorsConfig struct {
	Paths          []string `env:"CORS_PATHS,default=*"`
	AllowedMethods []string `env:"CORS_ALLOWED_METHODS,default=*"`
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS,default=*"`
	AllowedHeaders []string `env:"CORS_ALLOWED_HEADERS"`
	ExposedHeaders []string `env:"CORS_EXPOSED_HEADERS"`
	MaxAge         int      `env:"CORS_MAX_AGE,default=0"`
	SupportCred    bool     `env:"CORS_SUPPORT_CREDENTIALS,default=false"`
}

type Config interface {
	App() AppConfig
	Cors() CorsConfig
}

type ConfigImpl struct {
	app  AppConfig
	cors CorsConfig
}

func (c ConfigImpl) App() AppConfig {
	return c.app
}

func (c ConfigImpl) Cors() CorsConfig {
	return c.cors
}

func NewConfig(i *do.Injector) (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	var cfg AppConfig
	err := envdecode.StrictDecode(&cfg)
	if err != nil {
		return nil, err
	}
	var cors CorsConfig
	err = envdecode.StrictDecode(&cors)

	return &ConfigImpl{cfg, cors}, err
}

func init() {
	do.Provide[Config](do.DefaultInjector, NewConfig)
}
