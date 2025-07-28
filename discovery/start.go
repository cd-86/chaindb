package discovery

import (
	"time"
)

// INTERVAL: 网络环境中检查各类事项的间隔时间或超时时间.
func Start(interval time.Duration) {
	Register(interval)
	go func() {
		for {
			discovered_nodes := FindAll(interval)
			all_nodes := append(discovered_nodes, UserSpecifiedNodes.List()...)
			ActiveNodes.Store(&all_nodes)
		}
	}()
}
