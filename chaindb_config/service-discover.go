package chaindb_config

import "time"

// 每隔 多久 进行一次服务发现.
// 如果网络状况不怎么变化, 该值可以增大, 以降低服务发现的频率, 减少资源消耗.
const DiscoveryInterval = 10 * time.Second

// 使用 mDNS DNS-SD 进行服务发现的端口.
const DNSSDPort = 56782
