package config

import "bitbucket.org/dptsi/its-go/http"

func csrfConfig() http.CSRFConfig {
	return http.CSRFConfig{
		Methods: []string{"POST", "PUT", "PATCH", "DELETE"},
		Except:  []string{},
	}
}
