package blockchain

import (
	"container/list"
	"sync"
)

type Transaction struct {
	OwnerID uint32 `json:"OwnerID"`
	Nonce   uint32 `json:"Nonce"`
	Data    string `json:"Data"`
}

type TxQueue struct {
	transactions *list.List
	lock         sync.RWMutex
}

func (q *TxQueue) Enqueue(tx Transaction) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.transactions.PushBack(tx)
}

func (q *TxQueue) Remove(tx *list.Element) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.transactions.Remove(tx)
}
