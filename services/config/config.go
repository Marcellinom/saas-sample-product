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

type Config interface {
	App() AppConfig
}

type ConfigImpl struct {
	app AppConfig
}

func (c ConfigImpl) App() AppConfig {
	return c.app
}

func NewConfig(i *do.Injector) (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	var cfg AppConfig
	err := envdecode.StrictDecode(&cfg)

	return &ConfigImpl{cfg}, err
}

func init() {
	do.Provide[Config](do.DefaultInjector, NewConfig)
}
