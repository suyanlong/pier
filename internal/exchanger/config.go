package exchanger

import (
	"github.com/link33/sidecar/api"
	"github.com/link33/sidecar/internal/checker"
	"github.com/link33/sidecar/internal/peermgr"
	"github.com/link33/sidecar/internal/router"
	"github.com/link33/sidecar/internal/syncer"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/sirupsen/logrus"
)

type Config struct {
	checker   checker.Checker
	store     storage.Storage
	peerMgr   peermgr.PeerManager
	router    router.Router
	syncer    syncer.Syncer
	apiServer *api.Server
	logger    logrus.FieldLogger
}

type Option func(*Config)

func WithChecker(checker checker.Checker) Option {
	return func(config *Config) {
		config.checker = checker
	}
}

func WithPeerMgr(mgr peermgr.PeerManager) Option {
	return func(config *Config) {
		config.peerMgr = mgr
	}
}

func WithRouter(router router.Router) Option {
	return func(config *Config) {
		config.router = router
	}
}

func WithSyncer(syncer syncer.Syncer) Option {
	return func(config *Config) {
		config.syncer = syncer
	}
}

func WithAPIServer(apiServer *api.Server) Option {
	return func(config *Config) {
		config.apiServer = apiServer
	}
}

func WithStorage(store storage.Storage) Option {
	return func(config *Config) {
		config.store = store
	}
}

func WithLogger(logger logrus.FieldLogger) Option {
	return func(config *Config) {
		config.logger = logger
	}
}

func GenerateConfig(opts ...Option) *Config {
	config := &Config{}
	for _, opt := range opts {
		opt(config)
	}

	return config
}
