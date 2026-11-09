package redis

import "github.com/kelseyhightower/envconfig"

type config struct {
	Address  string `default:"localhost:6379" envconfig:"REDIS_ADDR"`
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
	DB       int    `default:"0" envconfig:"REDIS_DB"`
}

// newConfig loads Redis configuration values using the provided environment
// variable prefix.
func newConfig(envPrefix string) (*config, error) {
	var c config
	err := envconfig.Process(envPrefix, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
