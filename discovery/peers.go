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

// 无重复.
var AdministratorSpecifiedNodes = func() struct {
	Add    func(hosts []string)
	Remove func(hosts []string)
	List   func() []string
} {
	var (
		hosts []string
		lock  sync.RWMutex
	)

	Add := func(hosts_to_add []string) {
		lock.Lock()
		defer lock.Unlock()
		for _, host := range hosts_to_add {
			if slices.Contains(hosts, host) {
				continue
			}
			hosts = append(hosts, host)
		}
	}

	List := func() []string {
		lock.RLock()
		defer lock.RUnlock()
		return slices.Clone(hosts)
	}

	Remove := func(hosts_to_rm []string) {
		lock.Lock()
		defer lock.Unlock()
		hosts = slices.DeleteFunc(
			hosts,
			func(host string) bool {
				return slices.Contains(hosts_to_rm, host)
			},
		)
	}

	return struct {
		Add    func(hosts []string)
		Remove func(hosts []string)
		List   func() []string
	}{
		Add:    Add,
		Remove: Remove,
		List:   List,
	}
}()
