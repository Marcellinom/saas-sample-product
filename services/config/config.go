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

type AuthConfig struct {
	CookieDomain string `env:"AUTH_COOKIE_DOMAIN,default=localhost"`
	CookiePath   string `env:"AUTH_COOKIE_PATH,default=/"`
	MaxAge       int    `env:"AUTH_EXPIRATION,default=3600"`
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

type HTTPConfig struct {
	Port   int  `env:"HTTP_PORT,default=8080"`
	Secure bool `env:"HTTP_SECURE,default=false"`
}

type Config interface {
	App() AppConfig
	Auth() AuthConfig
	Cors() CorsConfig
	HTTP() HTTPConfig
}

type ConfigImpl struct {
	app  AppConfig
	auth AuthConfig
	cors CorsConfig
	http HTTPConfig
}

func (c ConfigImpl) App() AppConfig {
	return c.app
}

func (c ConfigImpl) Auth() AuthConfig {
	return c.auth
}

func (c ConfigImpl) Cors() CorsConfig {
	return c.cors
}

func (c ConfigImpl) HTTP() HTTPConfig {
	return c.http
}

func NewConfig(i *do.Injector) (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	var app AppConfig
	err := envdecode.StrictDecode(&app)
	if err != nil {
		return nil, err
	}
	var auth AuthConfig
	err = envdecode.StrictDecode(&auth)
	if err != nil {
		return nil, err
	}

	var cors CorsConfig
	err = envdecode.StrictDecode(&cors)
	if err != nil {
		return nil, err
	}

	var http HTTPConfig
	err = envdecode.StrictDecode(&http)

	return &ConfigImpl{app, auth, cors, http}, err
}

func init() {
	do.Provide[Config](do.DefaultInjector, NewConfig)
}
