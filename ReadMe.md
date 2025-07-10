# ChainDB: 基于区块链的 可信任网络下 键值存储型 内存数据库

## 介绍

### 目标

启动 ChainDB.
如果可以连接到其它设备, 加入到由那些设备组成的区块链网络中;
否则, 自己生成一个单一矿工节点的区块链网络.

ChainDB 进程退出后, 本地数据全部丢失.

⚠ 不建议 **在 ChainDB 运行期间** 切换 Wi-Fi,
可能会导致新 Wi-Fi 网络下的其它设备的数据被本设备的数据覆盖!

### 局限性

ChainDB 暂不支持回调监听.

## 概念

### 存储数据的最小单元

ChainDB 中存储一系列 transaction (事务), 每条 transaction 存储 `owner_uuid` 和 `data`.

同一个 owner 会有 0 个或多个 transaction 与其关联.
但 ChainDB 会记住 transaction 之间的先后顺序.

### 修改数据

区块链被认为是不可变的 (ChainDB 可能会主动丢弃某条 transaction, 但不会修改它),
但你仍然有办法 (在逻辑上) 修改 owner 所持有的 data.

例如, 对于 owner_uuid=42,
假设 ChainDB 中已经有一条与其关联的 transaction, 你可再插入一条 patch
`{"owner_uuid": 42, "data": "追加一个字符 'A'"}`, 同属于 owner_uuid=42 的一系列 data 可以按照
原始数据 + patches 的方式合并起来.

⚠ ChainDB 不关注 data 的具体内容, 如何组织 data 由用户自己决定!

### Confirmation Score

在分布式场景下, 一条 transaction 可能未能及时同步到其它所有主机, 或者最终因为不被区块链网络认可而被丢弃.
没有任何办法能确保一条 transaction 最终会被所有主机认可.

你可选择等待足够久的时间.
如果一条 transaction 一直存在于本地 ChainDB, 未被其主动丢弃,
那么随着时间的推移, 该 transaction 未被成功刻进区块链的概率是 **指数下降的**.

反映一条 transaction 同步程度的指标是 `confirmation_score`.
分数越高, 该 transaction 刻录失败的概率越低.

ChainDB 不提供指导, 用户应根据实际使用环境决定 `confirmation_score` 的最低阈值.
这是 比特币 交易的真实情况, ChainDB 遵循了 比特币 的设计.

## API

### 读

#### 列出 Owner

```
GetLocalOwners
```

返回 `Owners[{confirmation_score: <uint>, owner_uuid: <uint>}, ...]`:

```C++
struct /* Owner (匿名类) */ {
    std::uint64_t confirmation_score;
    std::uint64_t owner_uuid;
    // ... 其余的供内部使用的成员变量 ...
};
return std::vector</* Owner (匿名类) */>{
    {10, 0x65416},
    { 7, 0xFF464},
    { 7, 0x4154A},
    { 2, 0x1A4F6},
};  // 示例, 仅供参考.
```

#### 查询具体数据

```
GetLocalTxOwnedBy <owner_uuid>
```

返回 `Transactions[{confirmation_score: <uint>, data: "..."}, ...]`:

```C++
struct /* 事务类型 (匿名类) */ {
    std::uint64_t confirmation_score;
    std::string data;
    // ... 其余的供内部使用的成员变量 ...
};
return std::vector</* 事务类型 (匿名类) */>{
    {10, "Hello"},
    { 7, ", "},
    { 7, "world"},
    { 2,  "! "},
};  // 示例, 仅供参考.
```

### 写

```
AddTxLocally {
    required_confirmation_score: <uint>,
    transaction: {user_uuid: <uint>, data: "..."}
}
```

向本地 ChainDB 添加一条 transaction, 随后它可能会被整个区块链网络接受.

此处 `required_confirmation_score` 是一个用于控制同步时间的参数.
数值越大, 阻塞越久, 但 transaction 被刻进区块链的概率越大;
该值为 0 表示非阻塞调用, 你可将该值设为 0, 然后手动检查.

## 构建

### 配置

#### Confirmation 时间

设定大约多少秒增加一个 confirmation 数.
算法会根据设定值影响区块链的延伸速度.
