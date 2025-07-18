package discovery

import (
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

var ActiveNodes = func() *atomic.Pointer[[]net.IP] {
	var p atomic.Pointer[[]net.IP]
	p.Store(&[]net.IP{})
	return &p
}() // 不包括自己.

var UniqueNodeName = fmt.Sprintf("ChainDB-No%d", time.Now().UnixMilli())
