package app

import (
	"context"
	"fmt"
	_ "github.com/link33/sidercar/imports"
	"github.com/link33/sidercar/internal"
	"github.com/link33/sidercar/internal/appchain"
	"github.com/link33/sidercar/internal/loggers"
	"github.com/link33/sidercar/internal/manger"
	"github.com/link33/sidercar/internal/peermgr"
	"github.com/link33/sidercar/internal/repo"
	"github.com/link33/sidercar/internal/txcrypto"
	"github.com/link33/sidercar/pkg/plugins"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/meshplus/bitxhub-kit/storage/leveldb"
	"github.com/sirupsen/logrus"
	"path/filepath"
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

	manger *manger.Manager
}

// NewSidercar instantiates sidercar instance.
func NewSidercar(repoRoot string, config *repo.Config) (internal.Launcher, error) {
	store, err := leveldb.New(filepath.Join(config.RepoRoot, "store"))
	if err != nil {
		return nil, fmt.Errorf("read from datastaore %w", err)
	}

	logger := loggers.Logger(loggers.App)
	privateKey, err := repo.LoadPrivateKey(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("repo load key: %w", err)
	}

	addr, err := privateKey.PublicKey().Address()
	if err != nil {
		return nil, fmt.Errorf("get address from private key %w", err)
	}

	nodePrivKey, err := repo.LoadNodePrivateKey(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("repo load node key: %w", err)
	}

	var (
		//ck          checker.Checker
		//cryptor     txcrypto.Cryptor
		//ex internal.Launcher
		//sync        syncer.Syncer
		//apiServer   *api.Server
		peerManager peermgr.PeerManager
	)

	peerManager, err = peermgr.New(config, nodePrivKey, privateKey, 1, loggers.Logger(loggers.PeerMgr))
	if err != nil {
		return nil, fmt.Errorf("peerMgr create: %w", err)
	}

	peerManager.Ports()

	// sidercar网关客户端
	clients := plugins.CreateClients(config.Appchains, nil)

	var cryptor txcrypto.Cryptor
	//cryptor,err := txcrypto.NewDirectCryptor(peerManager,privateKey)
	//if err != nil {
	//	return nil, fmt.Errorf("txcrypto create: %w", err)
	//}

	clientPort := appchain.NewPorts(clients, cryptor, logger)

	ctx, cancel := context.WithCancel(context.Background())
	return &Sidercar{
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
		config: config,
	}, nil
}

// Start starts three main components of sidercar app
func (sidercar *Sidercar) Start() error {
	return nil
}

// Stop stops three main components of sidercar app
func (sidercar *Sidercar) Stop() error {
	return nil
}
