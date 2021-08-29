package app

import (
	"context"
	"github.com/meshplus/bitxhub-kit/crypto"
	"github.com/meshplus/bitxhub-kit/storage"
	"github.com/meshplus/pier/internal"
	"github.com/meshplus/pier/internal/executor"

	"github.com/hashicorp/go-plugin"
	"github.com/meshplus/pier/model/pb"
	"github.com/meshplus/pier/internal/lite"
	"github.com/meshplus/pier/internal/monitor"
	"github.com/meshplus/pier/internal/repo"
	"github.com/sirupsen/logrus"
)

// Pier represents the necessary data for starting the pier app
type Sidercar struct {
	privateKey crypto.PrivateKey
	//plugin     plugins.Client  //Client defines the interface that interacts with appchain 交互接口
	grpcPlugin *plugin.Client  //plugin 管理接口。可以定义到plugin里面
	monitor    monitor.Monitor //Monitor receives event from blockchain and sends it to network ：AppchainMonitor
	//Syncer 与 bithub交互的接口机制。
	exec      executor.Executor //represents the necessary data for executing interchain txs in appchain：与appchain链交互执行的接口
	lite      lite.Lite         //轻客户度
	storage   storage.Storage   // 存储
	exchanger internal.Launcher //主动交换，pier的核心动力引擎。
	ctx       context.Context
	cancel    context.CancelFunc
	//appchain   *appchainmgr.Appchain //appchain管理机制
	meta   *pb.Interchain //元数据
	config *repo.Config
	logger logrus.FieldLogger
}
