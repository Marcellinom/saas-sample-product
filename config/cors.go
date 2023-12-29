package config

import "bitbucket.org/dptsi/go-framework/http"

func corsConfig() http.CorsConfig {
	return http.CorsConfig{
		AllowedOrigins:   []string{"http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"content-type", "x-csrf-token"},
		ExposedHeaders:   []string{},
		MaxAge:           0,
		AllowCredentials: true,
	}
}
