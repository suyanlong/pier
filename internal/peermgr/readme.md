# peermgr
节点管理

放在每一个连接上，最低层。

router-》client-》peer
               -》plugin

封装plugin 与 peer 的上层关系。让router无法区分。
对router来说，它只知道，原地址、目标地址。
即封装成Port， 一个port代表是一个端点。
根据目标地址与原地址路由就可以了。
这个router里面还要体现出路由规则。
router是要能够获取路由规则，根据设置的规则路由、转发。
router 需要重写了。并且任务比较重要。

### 每个port也要与router交互
1、port转发消息到router
2、port接受消息到router
3、router根据接受消息的路由规则，即src|to地址进行转发、路由。
如果没有，则转发到上一层的sidercar。其实sidercar也是一个port。

### port 封装 pulgin与peer节点。
1、针对上层来说都是client。如何封装？主要还是根据链did标识或者sidercar did标识。
sidercar did标识是否要和自己绑定的链did标识一致、用同一个标识。
2、传输协议是否可以封装一下。一层一层的封装。
3、这就要求，不管是rpc、grpc、p2p都在下层。
4、因此peermgr实在p2p这个模块里面。
5、apiserver、router这些应该是在sidercar这个模块里面。
6、其它sidercar节点、pulgin都是client。
7、bithub连接点、其它链接点的都是plugin。
8、最终所有的都是port。只是对port进行细分，再细分。


## Monitor、Syncer这些接口，可以是每一个client里面应该实现的接口。
可以把每一个client的交易，发送、获取或者同步。只需要有标准化的一个数据协议（标准化协议）支持就可以了。

## moniter
Monitor receives event from blockchain and sends it to network


## exchanger(舍弃)
是整个核心上层，是一个整体循环，或者说系统的整个动力系统。
在sidercar里Router才是整个系统的动力系统。

## executor
ChannelExecutor represents the necessary data for executing interchain txs in appchain。
主要是和appchain交互，下沉到client下面，每个插件可以实现它。


## Checker规则验证
下沉到client层下面，或者是每一个用户侧，也可以是bithub一侧。


## rulemgr 规则验证管理，
主要是验证交易里面的存在性签名。
如果仅仅是转发的话。这个模块可以不用的。

sidercar 作用主要是路由转发，路由可以由用户设定。
路由规则优先级：
1、用户交易内部的路由规则最高。
2、用户在程序设定。
3、系统默认。

设计路由规则与用户验证规则结合。


## Syncer
主要是同步bitxhub里数据。
WrapperSyncer represents the necessary data for sync tx wrappers from bitxhub

放到插件里面去实现。可以作为每个input、output 里面client公共接口
或者，是每一个用户插件都应该实现的接口。

## Lite Client
轻客户端，主要是同步bithub上的数据。其实可以给每一个client使用。



## 整理

peer管理服务
节点同步需要有bithub同步过来。

提升路由规则到router最上层。
区分优先级

* 链ID 一条链可定是唯一的。
* 节点ID，每一个节点的ID都是唯一的（p2p ID）。
* 插件ID，每一个插件ID都是唯一的。插件绑定的链ID，才是主要路由的ID。
* sidercar ID每一个都是唯一的，可以用是p2p的ID充当。
* 


sidercar ID （如果目的地址（to）是自己，那就是转发或者存储。）
plugin ID
appchain ID 

如何路由？
plugin ID就是绑定的链ID（appchain ID）

sidercar ID 也可以作为路由的ID，绑定自身的唯一 peer ID，(用于指定路由)（转发消息）

other blockchain peer ID 这个就用链ID映射吧（直链）





































