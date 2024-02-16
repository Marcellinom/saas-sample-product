package config

import "github.com/dptsi/its-go/logging"

func loggingConfig() logging.Config {
	return logging.Config{
		Default: "go",
		Channels: map[string]logging.ChannelConfig{
			"go": {
				Driver:       "go",
				DriverConfig: nil,
			},
		},
	}
}
