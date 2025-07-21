package blockchain

import (
	"container/list"
	"sync"
)

type TxPool struct {
	transactions *list.List
	lock         sync.RWMutex
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: list.New(),
	}
}

func (q *TxPool) Enqueue(tx Transaction) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.transactions.PushBack(tx)
}

func (q *TxPool) Remove(tx *list.Element) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.transactions.Remove(tx)
}
