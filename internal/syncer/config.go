package syncer

import (
	rpcx "github.com/link33/sidercar/hub/client"
	"github.com/link33/sidercar/internal/repo"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/sirupsen/logrus"
)

type Config struct {
	client  rpcx.Client
	storage storage.Storage
	config  *repo.Config
	logger  logrus.FieldLogger
}

type Option func(*Config)

func WithClient(cli rpcx.Client) Option {
	return func(c *Config) {
		c.client = cli
	}
}

func WithStorage(store storage.Storage) Option {
	return func(config *Config) {
		config.storage = store
	}
}

func WithLogger(logger logrus.FieldLogger) Option {
	return func(config *Config) {
		config.logger = logger
	}
}

func GenerateConfig(opts ...Option) (*Config, error) {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}

	return config, nil
}
