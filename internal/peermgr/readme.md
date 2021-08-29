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
pier:在直连、联盟模式下是主要是和appchain交互，下沉到client下面，每个插件可以实现它。
pier:在中继模式下，使用创建BxhClient客户端代理。



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

* 定义转发
* 定义port
* 实现router
* 定义消息格式
* 定义各个接口

如果是这样，to、from才是正确的。


pier:
* UnionMode：pier之间组成联盟。路由与转发。
* 直链模式：pier 节点直链，不经过中继链。
* RelayMode：中继连架构。

## sidercar设计
1、跨链协议IBTPX设计：根据IBTP跨链协议增加路由策略：路由模式、路由方法等字段，一是用于验证交易的完整性，二是引入惩罚机制。
2、跨链网关sidercar架构设计，主要模块有：port模块、route模块、pulgin模块、monitor模块、exchanger模块、govern模块等。sidercar只做消息的转发，以及消息订阅，对消息不做任何的修改，存在作恶行为，任何其它节点都可以相互检举到治理模块，治理模块仲裁表决，裁掉作恶网关。
3、跨链网关路由接口port设计：抽象所有输入与输出为port接口，包括：blockchain peer、blockchain client、sidercar peer等全部抽象为port。
4、route模块：路由转发模块，根据IBTPX跨链扩展协议，使用路由策略，在port之间转发IBTPX消息。
5、pulgin模块：与其它区块链做跨链的各种适配插件。
6、govern模块：惩罚机制引入，主要根据IBTPX协议规定的路由路径，做节点校验，存在作恶行为，上报节点信息，相互监控。


* 不能完全照炒他们的东西。
* 常用的库可以使用。
* 定义的结构也可以使用。
* 设计自己简单的东西，才能完成跨链交互的东西。
* 功能进行必要的裁剪。设置实现优先级。
* 转发如何实现，丢失又该怎么办。
* 去掉rule、主从机制、加密机制先去掉。可以后面加进去。
* 


Syncer、Monitor、Executor三个都有相同的作用。是否可以重构为同一个接口。

Monitor、Executor:都是绑定自身的appchain的对象实现的接口，功能不一样。一个是读（监听）一个是写（提交）。
Syncer：是指中继架构下的hub客户端。
理解他们的代码逻辑，需要按照白皮书的设计来。

很多功能都不需要，只实现最简单的数据上链就可以了。他们实现的功能，很多都是白皮书写的东西，
许多的功能，这边可能都用不到。
比如实现最基本：
* 数据上链，（资产交换，后面再说）；
* 转发
* 路由
* 适配
这些是最基本的功能，后面的功能，可以持续迭代进去。

Monitor、Executor:都是绑定自身的appchain的对象实现的接口，一个pier绑定一条链ID，

如何做到别的网关收到以后，会转发给下一个网关，也就是说，自己这边没有appchainID,交由其它的网关进行处理。

不需要有恢复功能吧。

理解与读懂所有exchanger模块的实现。



## 输入与输出才有会ID，并且都是唯一。
* pluginID与绑定的blockchainID相同。
* pierID  各自唯一
* blockchainID 一条链的ID
* pluginID与blockchainID可能会相同。

* 注册路由表
* 删除ID
* 路由IBTPX数据包
* 校验
* 审计、治理
* pierID 做背书、签名、留言。
* 

## DID设计
#### 标识符格式
BitXHub 将运用数字身份来标识应用链、节点、链上账户甚至合约等各种实
体，其标识符格式略有差别：
其中，链的身份标识（Chain DID）的格式设计为：did:bitxhub:chain-name:.，
第一字段为 did 固定标识，第二字段为 BitXHub 网络的固定标识，第三字段为每
条链的链名，第四字段以.结尾。如 did:bitxhub:relaychain001:. 。
对于每条链上基于账户地址的账户数字身份标识（Account DID）其格式为：
did:bitxhub:chain-name:address，其中第四个字段为用户的账户地址。如：
did:bitxhub:relaychain001:0x12345678。


## 设计路由


## 裁剪功能
* 主备
* 模式选择
* 

## 
from、to：代表链ID

而sidercar ID 暂时未用，除非指定

只有路由功能、适配、转发。
重传就不需要了吧。

监听用户的交易

### 待定
* 同步数据。
* 两阶段提交，事务。
* 保存数据。
* 治理。
* 



### 

协议层分开

路由层修改为IBTPX协议

入参是接口、出参是可以是具体类型（也可以是接口，这个接口要包含外部的调用情况）。这样可以减少代码优化。

















