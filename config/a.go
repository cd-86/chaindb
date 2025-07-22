package chaindb

import "time"

// 近期内 整个网络 从概率上讲 平均多久生成一个区块.
const BlockTime = 1 * time.Second

// 根据前 N 个区块的平均生成时间来调整顺利产生下一个区块的概率,
// 以使得区块的平均生成时间接近 BlockTime.
const DifficultyAdjustmentWindow = 6

const HTTPPort = 56780

const GRPCPort = 56781

const DNSSDPort = 56782
