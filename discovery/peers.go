package discovery

import (
	"slices"
	"sync"
	"sync/atomic"
)

// 自动发现的节点的 host 与 用户指定的 host 列表.
// 定期刷新.
var ActiveNodes = func() *atomic.Pointer[[]string] {
	var p atomic.Pointer[[]string]
	p.Store(&[]string{})
	return &p
}() // 不包括自己.

var UserSpecifiedNodes = func() (
	nodes struct {
		hosts  []string
		lock   sync.RWMutex
		Add    func(host []string)
		Remove func(host []string)
		List   func() []string
	},
) {
	nodes.Add = func(hosts []string) {
		nodes.lock.Lock()
		defer nodes.lock.Unlock()
		nodes.hosts = append(nodes.hosts, hosts...)
	}
	nodes.Remove = func(hosts []string) {
		new_hosts_list := slices.DeleteFunc(
			slices.Clone(nodes.hosts),
			func(host string) bool {
				return slices.Contains(hosts, host)
			},
		)
		nodes.lock.Lock()
		defer nodes.lock.Unlock()
		nodes.hosts = new_hosts_list
	}
	nodes.List = func() []string {
		nodes.lock.RLock()
		defer nodes.lock.RUnlock()
		return slices.Clone(nodes.hosts)
	}
	return
}()
