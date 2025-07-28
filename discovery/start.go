package discovery

import (
	"slices"
	"time"
)

// INTERVAL: 网络环境中检查各类事项的间隔时间或超时时间.
func Start(interval time.Duration) {
	Register(interval)
	go func() {
		for {
			all_peers := slices.Concat(
				FindAll(interval),
				AdministratorSpecifiedNodes.List(),
			)
			ActiveNodes.Store(&all_peers)
		}
	}()
}
