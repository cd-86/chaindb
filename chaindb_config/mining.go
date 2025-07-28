package chaindb_config

import (
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

// 近期内 整个区块链网络 从概率上讲 平均多久生成一个区块.
// 生产环境下, 同一区块链网络下的所有节点都应使用相同的值.
const BlockTime = 1 * time.Second

// 用于唯一标识矿工.
// 该值可以自己指定, 当前实现是使用主机名和一个随机整数.
var MinerAddress = func() string {
	hostname, err := os.Hostname()
	short_hostname := strings.SplitN(hostname, ".", 2)[0]
	if err != nil {
		panic(err)
	}
	return short_hostname + "-" + strconv.Itoa(rand.Int())
}()
