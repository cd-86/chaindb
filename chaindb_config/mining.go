package chaindb_config

import (
	"math/rand/v2"
	"os"
	"strconv"
	"time"
)

// 近期内 整个区块链网络 从概率上讲 平均多久生成一个区块.
// 生产环境下, 同一区块链网络下的所有节点都应使用相同的值.
const BlockTime = 1 * time.Second

// 根据前 N 个区块的平均生成时间来调整顺利产生下一个区块的概率,
// 以使得区块的平均生成时间接近 BlockTime.
const DifficultyAdjustmentWindow = 6

// 用于唯一标识矿工.
// 该值可以自己指定, 当前实现是使用主机名和一个随机整数.
var MinerAddress = func() string {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return hostname + "-" + strconv.Itoa(rand.Int())
}()
