package api

import "github.com/kelseyhightower/envconfig"

// Config contains application settings loaded from environment variables.
type Config struct {
	ServiceName string `envconfig:"APP_NAME" default:"bookmark-service"`
	AppPort     int    `envconfig:"APP_PORT" default:"8081"`
	InstanceID  string `envconfig:"INSTANCE_ID" default:""`
}

// NewConfig loads application settings from environment variables.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
