package chaindb_config

// 矿工挖出新区块后会将其广播出去, 发送到其它节点的该端口.
// 不期待其它节点收到该 block, 整个区块链网络中没有任何矿工
// 接收到该 block 的概率不大.  就算真的一个都没有也无所谓.
const UDPPortPickBlock = 56782

// 分发 block 的 PCDN 端口, 使用 gRPC.
const TCPPortBlockPCDN = 56783
