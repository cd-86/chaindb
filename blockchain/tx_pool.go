package blockchain

import (
	"slices"
	"sync"
)

type TxPool struct {
	transactions []Transaction
	lock         sync.RWMutex
} // 按照 (OwnerID, Nonce) 升序.

// 如果 OwnerID 和 Nonce 相同, 则覆盖掉旧的.
func (pool *TxPool) Add(txs ...Transaction) {
	pool.lock.Lock()
	defer pool.lock.Unlock()

	for _, tx := range txs {
		pos, exist := slices.BinarySearchFunc(
			pool.transactions,
			tx,
			func(tx1, tx2 Transaction) int {
				if tx1.OwnerID != tx2.OwnerID {
					return int(tx1.OwnerID) - int(tx2.OwnerID)
				}
				return int(tx1.Nonce) - int(tx2.Nonce)
			},
		)

		if exist {
			pool.transactions[pos] = tx
		} else {
			pool.transactions = slices.Insert(pool.transactions, pos, tx)
		}
	}
}

func (pool *TxPool) PopFunc(txout chan<- Transaction, yes func(tx Transaction) bool) {
	txs := func() (txs []Transaction) {
		pool.lock.Lock()
		defer pool.lock.Unlock()

		txs, pool.transactions = pool.transactions, nil

		return
	}() // exchange

	for _, tx := range txs {
		if yes(tx) {
			txout <- tx
		} else {
			pool.Add(tx)
		}
	}
	close(txout)
}
