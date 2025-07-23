package discovery

import (
	"net"
	"sync/atomic"
)

var ActiveNodes = func() *atomic.Pointer[[]net.IP] {
	var p atomic.Pointer[[]net.IP]
	p.Store(&[]net.IP{})
	return &p
}() // 不包括自己.
