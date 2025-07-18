package discovery

import "testing"

func TestPing(t *testing.T) {
	ip_addresses_reachable := []string{
		"127.0.0.1",
		"localhost",
		"::1",
		// "bing.com",  有些环境下可能无法访问.
	}
	ip_addresses_unreachable := []string{
		"10.0.0.0", "10.255.255.255",
		"unknownhost",
		"WTF.shynur.fun",
	}

	for _, addr := range ip_addresses_reachable {
		err := ping(addr)
		if err != nil {
			t.Errorf("`ping %s' failed: %v", addr, err)
		}
	}

	for _, addr := range ip_addresses_unreachable {
		err := ping(addr)
		if err == nil {
			t.Errorf("`ping %s' should have failed, but it succeeded", addr)
		}
	}
}
