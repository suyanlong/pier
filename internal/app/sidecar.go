package app

import (
	"context"
	"path/filepath"

	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/meshplus/bitxhub-kit/storage/leveldb"
	"github.com/sirupsen/logrus"

	_ "github.com/link33/sidecar/imports"
	"github.com/link33/sidecar/internal"
	"github.com/link33/sidecar/internal/appchain"
	"github.com/link33/sidecar/internal/loggers"
	"github.com/link33/sidecar/internal/manger"
	"github.com/link33/sidecar/internal/peermgr"
	"github.com/link33/sidecar/internal/port"
	"github.com/link33/sidecar/internal/repo"
	"github.com/link33/sidecar/internal/txcrypto"
	"github.com/link33/sidecar/pkg/plugins"
)

// Sidecar represents the necessary data for starting the sidecar app
type Sidecar struct {
	privateKey crypto.PrivateKey
	storage    storage.Storage
	ctx        context.Context
	cancel     context.CancelFunc
	config     *repo.Config
	logger     logrus.FieldLogger
	manger     internal.Launcher
}

// NewSidecar instantiates sidecar instance.
func NewSidecar(repoRoot string, config *repo.Config) (internal.Launcher, error) {
	store, err := leveldb.New(filepath.Join(config.RepoRoot, "store"))
	Asset(err)
	logger := loggers.Logger(loggers.App)
	privateKey, err := repo.LoadPrivateKey(repoRoot)
	Asset(err)
	addr, err := privateKey.PublicKey().Address()
	Asset(err)
	nodePrivKey, err := repo.LoadNodePrivateKey(repoRoot)
	Asset(err)
	var (
	//ck          checker.Checker
	//cryptor     txcrypto.Cryptor
	//apiServer   *api.Server
	)
	portMap := port.NewPortMap()
	pm, err := peermgr.New(config, portMap, nodePrivKey, privateKey, 1, loggers.Logger(loggers.PeerMgr))
	Asset(err)
	clients := plugins.CreateClients(config.Appchains, nil)
	persister := manger.NewPersister(addr.String(), store, loggers.Logger(loggers.Manger))
	appchainMgr := appchainmgr.New(persister)
	cryptor, err := txcrypto.NewDirectCryptor(appchainMgr, privateKey)
	Asset(err)
	clientPort := appchain.NewPorts(clients, cryptor, logger)
	portMap.Adds(clientPort)
	mg, err := manger.NewManager(addr.String(), portMap, pm, appchainMgr, loggers.Logger(loggers.Manger))
	Asset(err)
	ctx, cancel := context.WithCancel(context.Background())
	return &Sidecar{
		storage: store,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
		config:  config,
		manger:  mg,
	}, nil
}

// Start starts three main components of sidecar app
func (s *Sidecar) Start() error {
	return s.manger.Start()
}

// Stop stops three main components of sidecar app
func (s *Sidecar) Stop() error {
	return s.manger.Stop()
}

func Asset(err error) {
	if err != nil {
		panic(err)
	}
}
