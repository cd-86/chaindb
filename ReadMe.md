# ChainDB: 基于区块链的可信任网络去中心化存储

## 目标

## 概念

ChainDB 中存储一系列 transaction (事务), 每条 transaction 存储 owner_id 和 data.  <br>
同一个 owner 会有 0 个或多个 transaction 与其关联.

区块链中的 transaction 被认为是不可变的 (ChainDB 可能会主动丢弃某条 transaction, 但不会修改它),
但你仍然有办法修改 owner 所持有的 data.  <br>
例如, 对于 owner_id=42, 假设 ChainDB 中已经有一条与其关联的 transaction, 你可再插入一条:
`{"owner_id": 42, "data": "追加一个字符 'A'"}`.



## API

查询: `GetLocalTxListOwnedBy <owner_id> `  <br>
返回 `Transactions[{confirmation_score: <uint>, data: "..."}, ...]`  <br>
(在区块链网络中, `confirmation_score` 表示一条记录 (`data`) 的可信程度.
 ChainDB 不提供指导, 调用方应根据实际使用环境决定 `confirmation_score` 的最低阈值.)



## 架构

