package app

import (
	"context"
	"path/filepath"

	appchainmgr "github.com/meshplus/bitxhub-core/appchain-mgr"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/meshplus/bitxhub-kit/storage/leveldb"
	"github.com/sirupsen/logrus"

	_ "github.com/link33/sidercar/imports"
	"github.com/link33/sidercar/internal"
	"github.com/link33/sidercar/internal/appchain"
	"github.com/link33/sidercar/internal/loggers"
	"github.com/link33/sidercar/internal/manger"
	"github.com/link33/sidercar/internal/peermgr"
	"github.com/link33/sidercar/internal/port"
	"github.com/link33/sidercar/internal/repo"
	"github.com/link33/sidercar/internal/txcrypto"
	"github.com/link33/sidercar/pkg/plugins"
)

// Sidercar represents the necessary data for starting the sidercar app
type Sidercar struct {
	privateKey crypto.PrivateKey
	//plugin     plugins.Client  //Client defines the interface that interacts with appchain 交互接口
	//grpcPlugin *plugin.Client  //plugin 管理接口。可以定义到plugin里面
	//monitor    monitor.Monitor //Monitor receives event from blockchain and sends it to network ：AppchainMonitor
	//Syncer 与 bithub交互的接口机制。
	//exec executor.Executor //represents the necessary data for executing interchain txs in appchain：与appchain链交互执行的接口
	storage storage.Storage // 存储
	//appchain   *appchainmgr.Appchain //appchain管理机制
	//exchanger internal.Launcher //主动交换，sidercar的核心动力引擎。
	ctx    context.Context
	cancel context.CancelFunc
	config *repo.Config
	logger logrus.FieldLogger
	manger internal.Launcher
}

// NewSidercar instantiates sidercar instance.
func NewSidercar(repoRoot string, config *repo.Config) (internal.Launcher, error) {
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
	//ex internal.Launcher
	//sync        syncer.Syncer
	//apiServer   *api.Server
	//peerManager peermgr.PeerManager
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
	return &Sidercar{
		storage: store,
		logger:  logger,
		ctx:     ctx,
		cancel:  cancel,
		config:  config,
		manger:  mg,
	}, nil
}

// Start starts three main components of sidercar app
func (s *Sidercar) Start() error {
	return s.manger.Start()
}

// Stop stops three main components of sidercar app
func (s *Sidercar) Stop() error {
	return s.manger.Stop()
}

func Asset(err error) {
	if err != nil {
		panic(err)
	}
}
