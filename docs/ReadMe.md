# ChainDB: 基于区块链的键值存储型数据库

<div align="center">

[![GoDoc](https://pkg.go.dev/badge/github.com/shynur/chaindb)](https://pkg.go.dev/github.com/shynur/chaindb)
![GoReportCard](https://goreportcard.com/badge/shynur/chaindb)
![测试覆盖率](https://codecov.io/gh/shynur/chaindb/graph/badge.svg)  <br />
[![可成功构建的徽章](https://github.com/shynur/chaindb/actions/workflows/go-build.yaml/badge.svg)](https://github.com/shynur/chaindb/actions/workflows/go-build.yaml)
[![通过测试的徽章](https://github.com/shynur/chaindb/actions/workflows/go-test.yaml/badge.svg)](https://github.com/shynur/chaindb/actions/workflows/go-test.yaml)  <br />
![代码行数统计徽章](https://tokei.rs/b1/github/shynur/chaindb?category=lines&label=代%20码%20行%20数&style=flat)  <!-- category=code 更真实 -->
![文件数统计徽章](https://tokei.rs/b1/github/shynur/chaindb?category=files&label=文%20件%20数%20目&style=flat)

</div>

## Features

基于区块链:
- 完全去中心化
- 动态规模: 任意且未知的节点可随时加入或离开
- 最终一致性
- 原子性修改
- 秒级可见性 (可配置)
- 单一节点可作为跳板连通两个网络

基于 DNS-SD:
- 自动发现

基于 Go 语言:
- 0 内存泄漏
- 高并发
- 跨平台 (GNU/Linux, MS-Windows)
- 低内存占用

基于 gRPC:
- 稳健的节点间数据同步

## 简介

### 目标

启动 ChainDB.
如果可以连接到其它设备, 加入到由那些设备组成的区块链网络中;
否则, 自己生成一个单一矿工节点的区块链网络.

> [!TIP]
> 如果设备 A 与设备 B 连接, 设备 B 与设备 C 连接, 则 A 与 C 可以间接连接, B 充当跳板.
> 跳板越少, 数据同步越慢, 且会使 *平均区块时间* (见后文) 变得不稳定.

只要加入到区块链网络, 即可同步数据.
因此现场可用笔记本电脑作为节点接入网络拉取数据库, 作为网络的接入点发送各类请求.

### 局限性

- 暂不支持回调监听.  (不影响效率, 区块链本来就是秒级更新, 用户哪怕每秒查询一次也是开销极小.)
- 暂不支持持久化存储.  (可在业务层实现.)

## 概念

### 存储数据的最小单元

概念上, ChainDB 是一个巨大的 `OwnerID:uint` $\mapsto$ `[history_1, history_2, ...]` 的 map.

ChainDB 中存储一系列 transactions, 每条 transaction 存储 `OwnerID` (key) 和 `Data` (value), 此外还有一个 `Nonce` 字段 (自然数).  <br />
Transaction 代表一条修订记录.
例如,

```json
{"OwnerID": 42, "Nonce": 0, "Data": "初始数据"}
{"OwnerID": 42, "Nonce": 1, "Data": "打个 patch"}
```

表示 `map[42] = ["初始数据", "打个 patch"]`.
其中 `Nonce` 表示第几次修订.

用户通过发送 transaction 向 ChainDB 中插入数据.  <br />
实践上:
- 如果用户需要新建一个 key-value pair, 他应该用 `{"OwnerID": <随机数, 此处假设是 42>, "Nonce": 0, "Data": <自定义>}`.
  - 用随机数作为 `OwnerID` 使得 key 基本不可能冲突.
  - 将 `Nonce` 置为 0 以解决潜在的冲突: 因为 ChainDB 保证同一个 owner 的一些系列 transactions 的 Nonce 从 0 开始按自然数依次递增.
    如果 key 有冲突 (即 `OwnerID=42` 已存在), 即 `{"OwnerID": 42, "Nonce": 0, "Data": "略"}` 已存在, 则本次 transaction 请求会被 ChainDB 丢弃.
- 如果需要向已有的 key (此处假设是 42) 追加 value, 应该用 `{"OwnerID": 42, "Nonce": map[42].length, "Data": <自定义>}`.
  - 将 `Nonce` 置为 `map[OwnerID].length` (此处假设是 1)
    - 用于标识 transaction 的唯一性以解决数据竞争: 因为 ChainDB 最终只保留一个合法的 transaction 请求, 使 `map[42][1]` 是一个标量.
    - 使 `Nonce` 自增.

### 区块链网络

#### 区块

主体部分是一个 transactions 数组.

此外, 它还有 `UUID` 和 `ParentUUID`, 使众多区块在逻辑上构成一个链表.

区块的 `Height` 表示它是自身所在的链表中的第几个区块.

#### 节点

每个运行 ChainDB 进程的设备都成为区块链网络里的一个节点.

节点负责将来自请求队列中的 transactions 打包成区块并追加到由本机维护的区块链末端.

用户可以通过 (例如, localhost 运行了 ChainDB)

```http
POST http://localhost:56784/api/v1/transactions?confirmation=0
Content-Type: application/json

{
    "OwnerID": 42,
    "Nonce": 0,
    "Data": "略"
}
```

请求将上述 transaction 打包进区块, 并追加到由 某个/某些 节点 (不一定是 localhost) 各自维护的区块链中.

#### 区块链

所有的链都 fork 自同一个链, 这些链此后又会 fork 出新的链.

```
              +--- [block-2]     +--- [block-5]
              |                  |     ^ 节点 A 的工作链, 一段时间后大概率会切换到 [block-9] 所属链
[block-1] <---+--- [block-3] <---+
              |                  |                    v 节点 B 的工作链, 即将 fork 出新链, 但可能被丢弃
              |                  +--- [block-6] <--- [block-7] <--- [block-8] <--- [block-9]
              +--- [block-4]                                                        ^ 节点 C 的工作链, [block-9] 即将被广播
```

> 最长链作为整个网络的共识链.

当节点构造出一个区块后, 会将其广播出去.

当一个节点接收到来自其它节点的区块后, 会根据该区块的 `Height` 来比较该区块所属链和本机所维护的链的长度.
**如果外来区块的所属链更长, 则将本机所维护的链切换到外来链的链尾并同步历史数据, 在此基础上向后延申, 原本的链被丢弃!**

> [!CAUTION]
> **不建议在 ChainDB 运行期间切换 Wi-Fi**,
> 否则可能会导致新 Wi-Fi 网络下所有节点的数据被本设备的数据覆盖.

因此, 网络中的多条链最终只有一条会成为 **主干**.  <br />
ChainDB 在概率上控制 **主干** 的增长速度为 1.x sec/block (试验中是 1.25), 即 *平均区块时间*.

> [!IMPORTANT]  
> 如果出于某些原因, 一个网络分裂成了两个, 那么这两个网络各自的 *平均区块时间* 仍然不变.
> ChainDB 将任何网络视为独立且平等的, 甚至是单一节点网络.

### Confirmation Score

ChainDB 实现了最终一致性.

`ConfirmationScore` 表示链上某个可查询到的 (说明该链此时尚未被丢弃) 区块后追加了几个新区块.
新区块越多, 说明该链 (或者叫 fork) 越有可能成为最长链.  <br />
一个区块的 `ConfirmationScore` 平均一个 *平均区块时间* 增加 1.

> [!TIP]
> 一个区块的 `ConfirmationScore=0` 说明它是某个 fork 的链尾.

Transaction 的 `ConfirmationScore` 就是该 transaction 所属区块的 `ConfirmationScore`.
`ConfirmationScore` 反映 transaction 所属链最终成为共识链的概率.

### 数学证明

#### 最终一致性

假设 A 维护的链的高度是 $H_a$, B 的是 $H_b$, 且

$$
z = H_a - H_b \gt 0
$$

- A 链每个 *平均区块时间* 生成一个区块;
- 令 $l = \\text{网络丢包率}$, 如果丢包则节点 B 不会切换到 A 链, 因为它没有收到 A 的区块广播.
  则 B 链每个 *平均区块时间* 有 $l$ 的概率生成一个区块.

当 B 链落后 A 链 $z$ 个区块时, B 链追上 A 链的概率

$$
p_z = l^z \in [0, 1)
$$

($l\neq1$, 否则这根本不是一个网络, 而是两个独立的节点.)

> 当 $l=50\\%$ 时, $p_{5}=3.125\\%$.

因此, 在一个良好的网络环境下, 共识链极难被回滚.

#### 数据可信度

假设用户向网络发送了一条 transaction request,
一段时间后, 从某个节点查询到该 transaction 的 `ConfirmationScore` 为 $c$.  <br />
可以计算该节点所维护的链属于 **主干** 的概率.

建模如下:

> 假设有 $N$ 个 counters, 初始值都是 $0$, 颜色各不相同.  <br />
> 每天中午随机选取一个 counter, 使其增加 $1$.
>
> 每天晚上, 查看当天自增的那个 counter 的颜色和数值.  <br />
> 对于每个数值小于那个 counter 的 counter, 有 $(1-l)$ 的概率被 那个 counter 的值和颜色覆盖.
>
> 然后开始下一天.
> 
> Question:
> 我们只查询特定的名为 T 的 counter.
> 某天早上观察到 T 的值为 $c$ 且它的颜色和自身的初始颜色相同 (可能从未改变, 也可能是被染回去了),
> 那么 T 和 数值最大的某个 counter 是同色的概率是多少?
