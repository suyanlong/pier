package app

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/hashicorp/go-plugin"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/meshplus/bitxhub-kit/storage/leveldb"
	"github.com/meshplus/bitxhub-model/pb"
	"github.com/meshplus/pier/api"
	rpcx "github.com/meshplus/pier/hub/client"
	_ "github.com/meshplus/pier/imports"
	"github.com/meshplus/pier/internal/agent"
	"github.com/meshplus/pier/internal/appchain"
	"github.com/meshplus/pier/internal/checker"
	"github.com/meshplus/pier/internal/exchanger"
	"github.com/meshplus/pier/internal/executor"
	"github.com/meshplus/pier/internal/lite"
	"github.com/meshplus/pier/internal/lite/direct_lite"
	"github.com/meshplus/pier/internal/lite/lite33"
	"github.com/meshplus/pier/internal/loggers"
	"github.com/meshplus/pier/internal/monitor"
	"github.com/meshplus/pier/internal/peermgr"
	"github.com/meshplus/pier/internal/repo"
	"github.com/meshplus/pier/internal/router"
	"github.com/meshplus/pier/internal/syncer"
	"github.com/meshplus/pier/internal/txcrypto"
	"github.com/meshplus/pier/pkg/plugins"
	_ "github.com/meshplus/pier/pkg/single"
	"github.com/sirupsen/logrus"
)

type Launcher interface {
	Start() error
	Stop() error
}

// Pier represents the necessary data for starting the pier app
type Pier struct {
	//privateKey crypto.PrivateKey
	//plugin     plugins.Client  //Client defines the interface that interacts with appchain 交互接口
	grpcPlugin *plugin.Client  //plugin 管理接口。可以定义到plugin里面
	monitor    monitor.Monitor //Monitor receives event from blockchain and sends it to network ：AppchainMonitor
	//Syncer 与 bithub交互的接口机制。
	exec      executor.Executor    //represents the necessary data for executing interchain txs in appchain：与appchain链交互执行的接口
	lite      lite.Lite            //轻客户度
	storage   storage.Storage      // 存储
	exchanger exchanger.IExchanger //主动交换，pier的核心动力引擎。
	ctx       context.Context
	cancel    context.CancelFunc
	//appchain   *appchainmgr.Appchain //appchain管理机制
	meta   *pb.Interchain //元数据
	config *repo.Config
	logger logrus.FieldLogger
}

// NewPier instantiates pier instance.
func NewPier(repoRoot string, config *repo.Config) (Launcher, error) {
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
		ck        checker.Checker
		cryptor   txcrypto.Cryptor
		ex        exchanger.IExchanger
		lite      lite.Lite
		sync      syncer.Syncer
		apiServer *api.Server
		meta      *pb.Interchain
		//chain       *appchainmgr.Appchain
		peerManager peermgr.PeerManager
	)

	switch config.Mode.Type {
	case repo.DirectMode: //直链模式
		peerManager, err = peermgr.New(config, nodePrivKey, privateKey, 1, loggers.Logger(loggers.PeerMgr))
		if err != nil {
			return nil, fmt.Errorf("peerMgr create: %w", err)
		}

		appchainMgr, err := appchain.NewManager(config.Appchain.DID, store, peerManager, loggers.Logger(loggers.AppchainMgr))
		if err != nil {
			return nil, fmt.Errorf("ruleMgr create: %w", err)
		}

		apiServer, err = api.NewServer(appchainMgr, peerManager, config, loggers.Logger(loggers.ApiServer))
		if err != nil {
			return nil, fmt.Errorf("gin service create: %w", err)
		}

		cryptor, err = txcrypto.NewDirectCryptor(appchainMgr, privateKey)
		if err != nil {
			return nil, fmt.Errorf("cryptor create: %w", err)
		}

		meta = &pb.Interchain{}
		lite = &direct_lite.MockLite{}
	case repo.RelayMode: //中继链架构模式。
		ck = &checker.MockChecker{}

		// pier register to bitxhub and got meta infos about its related
		// appchain from bitxhub
		opts := []rpcx.Option{
			rpcx.WithLogger(logger),
			rpcx.WithPrivateKey(privateKey),
		}
		nodesInfo := make([]*rpcx.NodeInfo, 0, len(config.Mode.Relay.Addrs))
		for _, addr := range config.Mode.Relay.Addrs {
			nodeInfo := &rpcx.NodeInfo{Addr: addr}
			if config.Security.EnableTLS {
				nodeInfo.CertPath = filepath.Join(config.RepoRoot, config.Security.Tlsca)
				nodeInfo.EnableTLS = config.Security.EnableTLS
				nodeInfo.CommonName = config.Security.CommonName
			}
			nodesInfo = append(nodesInfo, nodeInfo)
		}
		opts = append(opts, rpcx.WithNodesInfo(nodesInfo...), rpcx.WithTimeoutLimit(config.Mode.Relay.TimeoutLimit))
		client, err := rpcx.New(opts...)
		if err != nil {
			return nil, fmt.Errorf("create bitxhub client: %w", err)
		}

		// agent queries appchain info from bitxhub
		meta, err = getInterchainMeta(client, config.Appchain.DID)
		if err != nil {
			return nil, err
		}

		//chain, err = getAppchainInfo(client)
		//if err != nil {
		//	return nil, err
		//}

		//TODO 这里的中继是，bithub
		cryptor, err = txcrypto.NewRelayCryptor(client, privateKey)
		if err != nil {
			return nil, fmt.Errorf("cryptor create: %w", err)
		}

		lite, err = lite33.New(client, store, loggers.Logger(loggers.Lite33))
		if err != nil {
			return nil, fmt.Errorf("lite create: %w", err)
		}

		// 同步bithub数据。
		sync, err = syncer.New(addr.String(), config.Appchain.DID, repo.RelayMode,
			syncer.WithClient(client), syncer.WithLite(lite),
			syncer.WithStorage(store), syncer.WithLogger(loggers.Logger(loggers.Syncer)),
		)
		if err != nil {
			return nil, fmt.Errorf("syncer create: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported mode")
	}

	//use meta info to instantiate monitor and executor module
	extra, err := json.Marshal(meta.InterchainCounter)
	if err != nil {
		return nil, fmt.Errorf("marshal interchain meta: %w", err)
	}

	// pier网关客户端
	var cli plugins.Client
	var grpcPlugin *plugin.Client
	err = retry.Retry(func(attempt uint) error {
		cli, grpcPlugin, err = plugins.CreateClient(config.Appchain.DID, config.Appchain, extra)
		if err != nil {
			logger.Errorf("client plugin create:%s", err)
		}
		return err
	}, strategy.Wait(3*time.Second))
	if err != nil {
		logger.Panic(err)
	}

	// pier网关客户端
	mnt, err := monitor.New(cli, cryptor, loggers.Logger(loggers.Monitor))
	if err != nil {
		return nil, fmt.Errorf("monitor create: %w", err)
	}

	// 执行链一方。
	exec, err := executor.New(cli, config.Appchain.DID, store, cryptor, loggers.Logger(loggers.Executor))
	if err != nil {
		return nil, fmt.Errorf("executor create: %w", err)
	}

	ex, err = exchanger.New(config.Mode.Type, config.Appchain.DID, meta,
		exchanger.WithChecker(ck),
		exchanger.WithExecutor(exec),
		exchanger.WithMonitor(mnt),
		exchanger.WithPeerMgr(peerManager),
		exchanger.WithSyncer(sync),
		exchanger.WithAPIServer(apiServer),
		exchanger.WithStorage(store),
		exchanger.WithLogger(loggers.Logger(loggers.Exchanger)),
	)
	if err != nil {
		return nil, fmt.Errorf("exchanger create: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Pier{
		//privateKey: privateKey,
		//plugin:     cli,
		grpcPlugin: grpcPlugin,
		// appchain:   chain,
		meta:      meta,
		monitor:   mnt,
		exchanger: ex,
		exec:      exec,
		lite:      lite,
		logger:    logger,
		storage:   store,
		ctx:       ctx,
		cancel:    cancel,
		config:    config,
	}, nil
}

func NewUnionPier(repoRoot string, config *repo.Config) (Launcher, error) {
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
	var (
		ex          exchanger.IExchanger
		lite        lite.Lite
		sync        syncer.Syncer
		meta        *pb.Interchain
		ck          checker.Checker
		peerManager peermgr.PeerManager
	)

	ck = &checker.MockChecker{}
	nodePrivKey, err := repo.LoadNodePrivateKey(repoRoot)
	if err != nil {
		return nil, fmt.Errorf("repo load node key: %w", err)
	}

	// 连接指定节点。
	peerManager, err = peermgr.New(config, nodePrivKey, privateKey, config.Mode.Union.Providers, loggers.Logger(loggers.PeerMgr))
	if err != nil {
		return nil, fmt.Errorf("peerMgr create: %w", err)
	}

	opts := []rpcx.Option{
		rpcx.WithLogger(logger),
		rpcx.WithPrivateKey(privateKey),
	}
	// rpcx 连接。
	nodesInfo := make([]*rpcx.NodeInfo, 0, len(config.Mode.Union.Addrs))
	for _, addr := range config.Mode.Union.Addrs {
		nodeInfo := &rpcx.NodeInfo{Addr: addr}
		if config.Security.EnableTLS {
			nodeInfo.CertPath = filepath.Join(config.RepoRoot, config.Security.Tlsca)
			nodeInfo.EnableTLS = config.Security.EnableTLS
			nodeInfo.CommonName = config.Security.CommonName
		}
		nodesInfo = append(nodesInfo, nodeInfo)
	}
	opts = append(opts, rpcx.WithNodesInfo(nodesInfo...), rpcx.WithTimeoutLimit(config.Mode.Relay.TimeoutLimit))
	client, err := rpcx.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("create bitxhub client: %w", err)
	}

	meta = &pb.Interchain{}

	lite, err = lite33.New(client, store, loggers.Logger(loggers.Lite33))
	if err != nil {
		return nil, fmt.Errorf("lite create: %w", err)
	}

	sync, err = syncer.New(addr.String(), config.Appchain.DID, repo.UnionMode,
		syncer.WithClient(client), syncer.WithLite(lite), syncer.WithStorage(store), syncer.WithLogger(loggers.Logger(loggers.Syncer)),
	)
	if err != nil {
		return nil, fmt.Errorf("syncer create: %w", err)
	}

	cli := agent.CreateClient(client) //创建BxhClient客户端代理
	exec, err := executor.New(cli, addr.String(), store, nil, loggers.Logger(loggers.Executor))
	if err != nil {
		return nil, fmt.Errorf("executor create: %w", err)
	}

	r := router.New(peerManager, store, loggers.Logger(loggers.Router), peerManager.(*peermgr.Swarm).ConnectedPeerIDs())

	ex, err = exchanger.New(config.Mode.Type, addr.String(), meta,
		exchanger.WithExecutor(exec),
		exchanger.WithPeerMgr(peerManager),
		exchanger.WithChecker(ck),
		exchanger.WithSyncer(sync),
		exchanger.WithStorage(store),
		exchanger.WithRouter(r),
		exchanger.WithLogger(loggers.Logger(loggers.Exchanger)),
	)
	if err != nil {
		return nil, fmt.Errorf("exchanger create: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Pier{
		//privateKey: privateKey,
		meta:      meta,
		exchanger: ex,
		exec:      exec,
		lite:      lite,
		storage:   store,
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
		config:    config,
	}, nil
}

// Start starts three main components of pier app
func (pier *Pier) Start() error {
	if err := pier.lite.Start(); err != nil {
		pier.logger.Errorf("lite start: %w", err)
		return err
	}
	pier.logger.WithFields(logrus.Fields{
		"id":                     pier.meta.ID,
		"interchain_counter":     pier.meta.InterchainCounter,
		"receipt_counter":        pier.meta.ReceiptCounter,
		"source_receipt_counter": pier.meta.SourceReceiptCounter,
	}).Info("Pier information")
	if err := pier.monitor.Start(); err != nil {
		pier.logger.Errorf("monitor start: %w", err)
		return err
	}
	if err := pier.exec.Start(); err != nil {
		pier.logger.Errorf("executor start: %w", err)
		return err
	}
	err := pier.exchanger.Start()
	return err
}

// Stop stops three main components of pier app
func (pier *Pier) Stop() error {
	if err := pier.monitor.Stop(); err != nil {
		return fmt.Errorf("monitor stop: %w", err)
	}

	if err := pier.exec.Stop(); err != nil {
		return fmt.Errorf("executor stop: %w", err)
	}

	if err := pier.lite.Stop(); err != nil {
		return fmt.Errorf("lite stop: %w", err)
	}

	if err := pier.exchanger.Stop(); err != nil {
		return fmt.Errorf("exchanger stop: %w", err)
	}
	// stop appchain plugin first and kill plugin process
	pier.grpcPlugin.Kill()
	return nil
}
