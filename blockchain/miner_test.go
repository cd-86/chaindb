package blockchain

import (
	"math/rand/v2"
	"testing"
	"time"
)

func TestStartMining(t *testing.T) {
	local_chain := New()
	local_tx_pool := TxPool{}

	StartMining(local_chain, &local_tx_pool)

	stop := make(chan int)
	go func(stop <-chan int) {
		for i := 0; i != 1_0000; i++ {
			tx := Transaction{
				OwnerID: rand.Uint32N(3),
				Nonce:   rand.Uint32N(5),
			}
			local_tx_pool.Add(tx)

			time.Sleep(10 * time.Millisecond)

			select {
			case <-stop:
				return
			default:
			}
		}
	}(stop)
	time.Sleep(10 * time.Second)
	stop <- 0

	t.Logf("%+v", local_chain)
}
