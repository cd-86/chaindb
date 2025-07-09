# ChainDB: 基于区块链的 可信任网络下 键值存储型 内存数据库

## 目标

## 概念

### 存储数据的最小单元

ChainDB 中存储一系列 transaction (事务), 每条 transaction 存储 owner_id 和 data.

同一个 owner 会有 0 个或多个 transaction 与其关联.  <br>
但 ChainDB 会记住 transaction 之间的先后顺序.

### 修改数据

区块链被认为是不可变的 (ChainDB 可能会主动丢弃某条 transaction, 但不会修改它),
但你仍然有办法 (在逻辑上) 修改 owner 所持有的 data.

例如, 对于 owner_id=42,
假设 ChainDB 中已经有一条与其关联的 transaction, 你可再插入一条 patch
`{"owner_id": 42, "data": "追加一个字符 'A'"}`, 同属于 owner_id=42 的一系列 data 可以按照
原始数据 + patches 的方式合并起来.

⚠ ChainDB 不关注 data 的具体内容, 如何组织 data 由用户自己决定!

## API

查询: `GetLocalTxListOwnedBy <owner_id> `  <br>
返回 `Transactions[{confirmation_score: <uint>, data: "..."}, ...]`  <br>
(在区块链网络中, `confirmation_score` 表示一条记录 (`data`) 的可信程度.
 ChainDB 不提供指导, 调用方应根据实际使用环境决定 `confirmation_score` 的最低阈值.)



## 架构

