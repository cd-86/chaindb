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

type TxPool struct {
	transactions *list.List
	lock         sync.RWMutex
}

var LocalTxPool = TxPool{
	transactions: list.New(),
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
