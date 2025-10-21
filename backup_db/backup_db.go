package backup_db

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/shynur/chaindb/blockchain"
)

type BackupDB struct {
	db *sql.DB
}

func New(file string) (*BackupDB, error) {
	db, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}
	return &BackupDB{db: db}, nil
}

func (b *BackupDB) Close() error {
	return b.db.Close()
}

// 获取数据库最后一个区块的 UUID
func (b *BackupDB) GetLastBlockUUID() (uint64, error) {
	var uuid uint64
	err := b.db.QueryRow(`SELECT UUID FROM backup ORDER BY ID DESC LIMIT 1;`).Scan(&uuid)
	if err != nil {
		return 0, err
	}
	return uuid, nil
}

func (b *BackupDB) WriteBlocksToDB(block []blockchain.Block) error {
	// 是否有 backup 表， 如果没有则创建 表包含 Timestamp, MinerAddress, UUID, Height, ParentUUID  Transactions 字段
	_, err := b.db.Exec(`CREATE TABLE IF NOT EXISTS backup (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Timestamp REAL,
		MinerAddress TEXT,
		UUID TEXT,
		Height TEXT,
		ParentUUID TEXT,
		Transactions TEXT
	);`)
	if err != nil {
		return err
	}
	tx, err := b.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO backup (Timestamp, MinerAddress, UUID, Height, ParentUUID, Transactions) VALUES (?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, block := range block {
		transactionsBytes, err := json.Marshal(block.Transactions)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(block.Timestamp, block.MinerAddress, strconv.FormatUint(block.UUID, 10), strconv.FormatUint(uint64(block.Height), 10), strconv.FormatUint(block.ParentUUID, 10), string(transactionsBytes))
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (b *BackupDB) ReadBlocksFromDB() ([]blockchain.Block, error) {
	// 是否有 backup 表， 如果没有则创建 表包含 Timestamp, MinerAddress, UUID, Height, ParentUUID  Transactions 字段
	_, err := b.db.Exec(`CREATE TABLE IF NOT EXISTS backup (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		Timestamp REAL,
		MinerAddress TEXT,
		UUID TEXT,
		Height TEXT,
		ParentUUID TEXT,
		Transactions TEXT
	);`)
	if err != nil {
		return nil, err
	}
	rows, err := b.db.Query(`SELECT Timestamp, MinerAddress, UUID, Height, ParentUUID, Transactions FROM backup;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []blockchain.Block
	for rows.Next() {
		var block blockchain.Block
		var transactionsStr string
		err := rows.Scan(&block.Timestamp, &block.MinerAddress, &block.UUID, &block.Height, &block.ParentUUID, &transactionsStr)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(transactionsStr), &block.Transactions)
		blocks = append(blocks, block)
	}
	return blocks, nil
}
