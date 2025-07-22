package chaindb

import "time"

// 近期内 整个区块链网络 从概率上讲 平均多久生成一个区块.
// 生产环境下, 同一区块链网络下的所有节点都应使用相同的值.
const BlockTime = 1 * time.Second

// 根据前 N 个区块的平均生成时间来调整顺利产生下一个区块的概率,
// 以使得区块的平均生成时间接近 BlockTime.
const DifficultyAdjustmentWindow = 6

const HTTPPort = 56780

const GRPCPort = 56781

const DNSSDPort = 56782

// 每隔 多久 进行一次服务发现.
// 如果网络状况不怎么变化, 该值可以增大, 以降低服务发现的频率, 减少资源消耗.
const DiscoveryInterval = 10 * time.Second
