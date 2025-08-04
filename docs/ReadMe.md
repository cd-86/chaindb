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

- 完全去中心化
- 自动发现
- 规模动态变化, 任意未知节点可随时加入或离开
- 最终一致性
- 原子性修改
- 秒级可见性
- 0 内存泄漏
- 高并发
- 跨平台 (GNU/Linux, MS-Windows)
- 单一节点可作为跳板联通两个网络
- 低内存占用

## 介绍

### 目标

启动 ChainDB.
如果可以连接到其它设备, 加入到由那些设备组成的区块链网络中;
否则, 自己生成一个单一矿工节点的区块链网络.

> [!TIP]
> 如果设备 A 与设备 B 连接, 设备 B 与设备 C 连接, 则 A 与 C 可以间接连接, B 充当跳板.
> 跳板越少, 数据同步越慢, 且会影响 *平均区块时间*.

ChainDB 进程退出后, 本地数据全部丢失.
只要加入到区块链网络, 即可同步数据, 因此现场可用笔记本电脑 (可以是 MS-Windows, ChainDB 跨平台) 作为节点接入网络定期拉取数据库.

⚠ **不建议在 ChainDB 运行期间切换 Wi-Fi**,
可能会导致新 Wi-Fi 网络下的其它设备的数据被本设备的数据覆盖.

### 局限性

- 暂不支持回调监听.
- 暂不支持持久化存储 (可在业务层实现).

## 概念

### 存储数据的最小单元

ChainDB 中存储一系列 transaction (事务), 每条 transaction 存储 `OwnerID` (key) 和 `Data` (value).

同一个 owner 会有 0 个或多个 transaction 与其关联.
但 ChainDB 会记住 transaction 之间的先后顺序.

### 修改数据

区块链被认为是不可变的, 只能追加数据.
逻辑上讲, 追加数据可以理解为修改数据的一种方式.

例如, 对于 `OwnerID=42`,
假设 ChainDB 中已经有一条与其关联的 transaction, 你可继续插入
`{"OwnerID": 42, "Data": "追加一个字符 'A'"}`.

当查询 `OwnerID=42` 时, 会按追加顺序返回所有 transaction.

### Confirmation Score

在分布式场景下, 一条 transaction 可能未能及时同步到其它所有主机, 或者最终因为不被区块链网络认可而被丢弃.
没有任何办法能确保一条 transaction 最终会被所有主机认可.

你可选择等待足够久的时间.
如果一条 transaction 一直存在于本地 ChainDB, 未被其主动丢弃,
那么随着时间的推移, 该 transaction 未被成功刻进区块链的概率是 **指数下降的**.

反映一条 transaction 同步程度的指标是 `ConfirmationScore`.
分数越高, 该 transaction 刻录失败的概率越低.

ChainDB 不提供指导, 用户应根据实际使用环境决定 `ConfirmationScore` 的最低阈值.
这是 比特币 交易的真实情况, ChainDB 遵循了 比特币 的设计.
