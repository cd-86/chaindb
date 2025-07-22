package blockchain

import (
	"container/list"
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

func (q *TxPool) Clear_sync(tx *list.Element) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.transactions = nil
}
