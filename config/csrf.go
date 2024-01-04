package config

import "github.com/dptsi/its-go/http"

func csrfConfig() http.CSRFConfig {
	return http.CSRFConfig{
		Methods: []string{"POST", "PUT", "PATCH", "DELETE"},
		Except:  []string{},
	}
}
