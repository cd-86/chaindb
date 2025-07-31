package discovery

import (
	"log"
	"time"
)

// INTERVAL: 网络环境中检查各类事项的间隔时间或超时时间.
func Start(interval time.Duration) {
	Register(interval)
	go func() {
		for {
			all_peers := FindAllActive(interval)
			ActiveNodes.Store(&all_peers)
			log.Printf(
				"[ DNS-SD] 当前活跃节点列表:\n%s\n",
				func() (list_peers string) {
					for _, peer := range all_peers {
						list_peers += "\t" + peer + "\n"
					}
					return
				}(),
			)
		}
	}()
}
