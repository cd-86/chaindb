package blockchain

import (
	"sync"
)

type TxPool struct {
	transactions []Transaction
	lock         sync.RWMutex
}

func (q *TxPool) Add_sync(tx Transaction) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.transactions = append(q.transactions, tx)
}

func (q *TxPool) ExchangeNil_sync() (take []Transaction) {
	q.lock.Lock()
	defer q.lock.Unlock()
	take, q.transactions = q.transactions, nil
	return
}
