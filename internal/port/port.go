package port

import (
	"github.com/link33/sidercar/model/pb"
	"math/rand"
	"sync"
)

// port 类型：主要是sider peer、plugin、blockchain peer。
type Type int

const (
	Hub      = "hub"      //Hub: 同步数据，同步元数据等。
	Sidercar = "sidercar" //SiderCar节点
	Appchain = "appchain" //区块链客户端
)

// 设计一套port管理机制：包括各种的管理模块。以组合的行驶。
// 设计一套，管理机制。
// 与中继交互的是单独完整的机制。并且注册到路由表中。或者更加类型，这样就限制一个sidercar最多只能连接一个hub。避免网络风暴。或者只是一个转发功能。转发到指定节点。
// 先是从转发开始完成。
// 协议实现
// 路由策略
// 客户度
// 验证器（上链）
// 动力核心引擎：主动发起、主动查询、主动关联数据、存储数据。
// 应答的数据，就原路返回，根据请求路径，原路返回。
// 与hub交互一端的数据，需要记录数据。而与appchain 交互的数据也需要存储下来。

// 代表每一个路由端点
// 对router来说，只需要体现两个作用：1、did 唯一标识，2：接受一个ibtp数据函数。Send、recive（不对外开放）
// client 可是代表是sdk rpc 这些东西。
// 是否要做一个管理层，管理整个port.以及plugin。
// Port
type Port interface {
	ID() string
	Type() string
	Name() string
	Tag() string

	Send(msg Message) (*pb.Message, error)
	AsyncSend(msg Message) error
	ListenIBTPX() <-chan *pb.IBTPX
}

type Message interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

//获取唯一hub
//根据类型获取port
//根据ID获取
type PortMap struct {
	rw           sync.RWMutex
	peerPort     map[string]Port
	appchainPort map[string]Port
	hubPort      Port
}

func NewPortMap() *PortMap {
	return &PortMap{
		rw:           sync.RWMutex{},
		peerPort:     map[string]Port{},
		appchainPort: map[string]Port{},
		hubPort:      nil,
	}
}

func (p *PortMap) Adds(pp []Port) {
	p.rw.Lock()
	p.rw.Unlock()
	for _, pt := range pp {
		p.add(pt)
	}
}

func (p *PortMap) Add(pt Port) {
	p.rw.Lock()
	p.rw.Unlock()
	p.add(pt)
}

func (p *PortMap) add(pt Port) {
	switch pt.Type() {
	case Hub:
		if p.hubPort == nil {
			p.hubPort = pt
		}
	case Appchain:
		p.appchainPort[pt.ID()] = pt
	case Sidercar:
		p.peerPort[pt.ID()] = pt
	}
}

func (p *PortMap) getHub() (Port, bool) {
	p.rw.RLocker()
	defer p.rw.RUnlock()
	if p.hubPort == nil {
		return nil, false
	}
	return p.hubPort, true
}

func (p *PortMap) Port(id string) (Port, bool) {
	p.rw.RLocker()
	defer p.rw.RUnlock()
	if p.hubPort.ID() == id {
		return p.hubPort, true
	}
	if pt, is := p.peerPort[id]; is {
		return pt, is
	}
	if pt, is := p.appchainPort[id]; is {
		return pt, is
	}
	return nil, false
}

func (p *PortMap) RouterPortByID(ids []string) []Port {
	p.rw.RLocker()
	defer p.rw.RUnlock()
	var ports []Port
	for _, id := range ids {
		if pt, is := p.peerPort[id]; is {
			ports = append(ports, pt)
		}
	}
	return ports
}

func (p *PortMap) RouterPortByTag(tag string) []Port {
	p.rw.RLocker()
	defer p.rw.RUnlock()
	var ports []Port
	for _, pt := range p.peerPort {
		if tag == pt.Tag() {
			ports = append(ports, pt)
		}
	}
	return ports
}

func (p *PortMap) AllRouterPort() []Port {
	p.rw.RLocker()
	defer p.rw.RUnlock()
	var ports []Port
	for _, pt := range p.peerPort {
		ports = append(ports, pt)
	}
	return ports
}

func (p *PortMap) RandRouterPort() Port {
	p.rw.RLocker()
	defer p.rw.RUnlock()
	var randPort Port
	l := len(p.peerPort)
	i := rand.Intn(l + 1)
	j := 0
	for _, pt := range p.peerPort {
		j++
		if i == j {
			randPort = pt
		}
	}
	return randPort
}

func (p *PortMap) Remove(pt Port) {
	p.rw.Lock()
	p.rw.Unlock()
	p.remove(pt)
}

func (p *PortMap) remove(pt Port) {
	switch pt.Type() {
	case Hub:
		if p.hubPort.ID() == pt.ID() {
			p.hubPort = nil
		}
	case Appchain:
		delete(p.appchainPort, pt.ID())
	case Sidercar:
		delete(p.peerPort, pt.ID())
	}
}

func (p *PortMap) Removes(ppt []Port) {
	p.rw.Lock()
	p.rw.Unlock()
	for _, pt := range ppt {
		p.remove(pt)
	}
}

func (p *PortMap) Store(id string, port Port) {
	p.Add(port)
}

func (p *PortMap) Load(key string) (value Port, ok bool) {
	return p.Port(key)
}
