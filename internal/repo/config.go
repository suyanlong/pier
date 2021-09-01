package repo

import (
	"fmt"
	"github.com/meshplus/bitxhub-kit/fileutil"
	"github.com/meshplus/bitxhub-kit/types"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
	"time"
)

const (
	DirectMode = "direct" //直连模式
	RelayMode  = "relay"  //
	UnionMode  = "union"  //联盟模式，中继架构。代表的是部署架构，不够灵活。不能随时变动。
)

// Config represents the necessary config data for starting sidecar
type Config struct {
	RepoRoot  string
	Title     string    `toml:"title" json:"title"`
	Port      Port      `toml:"port" json:"port"`
	Log       Log       `toml:"log" json:"log"`
	Appchains Appchains `toml:"appchains" json:"appchains"`
	Security  Security  `toml:"security" json:"security"`
	Peer      Peer      `toml:"peer" json:"peer"`
}

type Peer struct {
	Peers      []string `toml:"peers" json:"peers"`
	Connectors []string `toml:"connectors" json:"connectors"`
	Providers  uint64   `toml:"providers" json:"providers"`
}

// Security are certs used to setup connection with tls
type Security struct {
	EnableTLS  bool   `mapstructure:"enable_tls"`
	Tlsca      string `toml:"tlsca" json:"tlsca"`
	CommonName string `mapstructure:"common_name" json:"common_name"`
}

// Port are ports providing http and pprof service
type Port struct {
	Http  int64 `toml:"http" json:"http"`
	PProf int64 `toml:"pprof" json:"pprof"`
}

// Log are config about log
type Log struct {
	Dir          string    `toml:"dir" json:"dir"`
	Filename     string    `toml:"filename" json:"filename"`
	ReportCaller bool      `mapstructure:"report_caller"`
	Level        string    `toml:"level" json:"level"`
	Module       LogModule `toml:"module" json:"module"`
}

type LogModule struct {
	ApiServer   string `mapstructure:"api_server" toml:"api_server" json:"api_server"`
	AppchainMgr string `mapstructure:"appchain_mgr" toml:"appchain_mgr" json:"appchain_mgr"`
	Lite33      string `mapstructure:"lite33" toml:"lite33" json:"lite33"`
	Exchanger   string `toml:"exchanger" json:"exchanger"`
	Executor    string `toml:"executor" json:"executor"`
	Monitor     string `toml:"monitor" json:"monitor"`
	PeerMgr     string `mapstructure:"peer_mgr" toml:"peer_mgr" json:"peer_mgr"`
	Router      string `toml:"router" json:"router"`
	Swarm       string `toml:"swarm" json:"swarm"`
	Syncer      string `toml:"syncer" json:"syncer"`
	Manger      string `toml:"manger" json:"manger"`
}

// Appchain are configs about appchain
type Appchain struct {
	Enable bool   `toml:"enable" json:"enable"`
	Type   string `toml:"type" json:"type"`
	DID    string `toml:"did" json:"did"`
	Config string `toml:"config" json:"config"`
	Plugin string `toml:"plugin" json:"plugin"`
}

type Appchains struct {
	Appchains []Appchain
}

// DefaultConfig returns config with default value
func DefaultConfig() *Config {
	return &Config{
		RepoRoot: "sidecar",
		Title:    "sidecar configuration file",
		Port: Port{
			Http:  8080,
			PProf: 44555,
		},
		Log: Log{
			Dir:          "logs",
			Filename:     "sidecar.log",
			ReportCaller: false,
			Level:        "info",
			Module: LogModule{
				AppchainMgr: "info",
				Exchanger:   "info",
				Executor:    "info",
				Lite33:      "info",
				Monitor:     "info",
				Swarm:       "info",
				Syncer:      "info",
				PeerMgr:     "info",
				Router:      "info",
				ApiServer:   "info",
			},
		},
		Peer: Peer{
			Peers:      []string{"localhost:60011", "localhost:60012", "localhost:60013", "localhost:60014"},
			Connectors: []string{},
			Providers:  1,
		},
		Appchains: Appchains{[]Appchain{
			{
				Enable: true,
				Type:   "appchain",
				DID:    "did:bitxhub:appchain:.",
				Plugin: "appchain_plugin",
				Config: "fabric",
			},
			{
				Enable: true,
				Type:   "hub",
				DID:    "did:bitxhub:chain33:.",
				Plugin: "appchain_plugin",
				Config: "chain33",
			},
		}},
	}
}

// UnmarshalConfig read from config files under config path
func UnmarshalConfig(repoRoot string) (*Config, error) {
	configPath := filepath.Join(repoRoot, ConfigName)

	if !fileutil.Exist(configPath) {
		return nil, fmt.Errorf("file %s doesn't exist, please initialize sidecar firstly", configPath)
	}

	viper.SetConfigFile(configPath)
	viper.SetConfigType("toml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SIDECAR")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := DefaultConfig()

	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	config.RepoRoot = repoRoot

	return config, nil
}

type Mode struct {
	Type  string `toml:"type" json:"type"`
	Relay Relay  `toml:"relay" json:"relay"`
	// TODO 连接节点
	Peers      []string `toml:"peers" json:"peers"`
	Connectors []string `toml:"connectors" json:"connectors"`
	Providers  uint64   `toml:"providers" json:"providers"`
}

// Relay are configs about bitxhub
type Relay struct {
	Addrs        []string      `toml:"addrs" json:"addrs"`
	TimeoutLimit time.Duration `mapstructure:"timeout_limit" json:"timeout_limit"`
	Quorum       uint64        `toml:"quorum" json:"quorum"`
	Validators   []string      `toml:"validators" json:"validators"`
}

type Direct struct {
	Peers []string `toml:"peers" json:"peers"`
}

type Union struct {
	Addrs      []string `toml:"addrs" json:"addrs"`
	Connectors []string `toml:"connectors" json:"connectors"`
	Providers  uint64   `toml:"providers" json:"providers"`
}

// GetValidators gets validator address of bitxhub
func (relay *Relay) GetValidators() []*types.Address {
	validators := make([]*types.Address, 0)
	for _, v := range relay.Validators {
		validators = append(validators, types.NewAddressByStr(v))
	}
	return validators
}
