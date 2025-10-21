package main

import (
	"log"
	"slices"
	"time"

	"github.com/shynur/chaindb/backup_db"
	"github.com/shynur/chaindb/blockchain"
	"github.com/shynur/chaindb/chaindb_config"
	"github.com/shynur/chaindb/discovery"
	"github.com/shynur/chaindb/monitor"
	"github.com/shynur/chaindb/validator"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := backup_db.New(chaindb_config.BackupDBFilePath)
	if err == nil {
		var blocks []blockchain.Block
		blocks, err = db.ReadBlocksFromDB()
		if err == nil {
			if len(blocks) == 0 {
				log.Printf("[ Backup] 备份数据库中没有区块可供恢复\n")
			} else {
				for _, block := range blocks {
					blockchain.BlockCache.Store(block.UUID, block)
				}
				if LocalChain.TrySwitchHead(blocks[len(blocks)-1].UUID) {
					log.Printf("[ Backup] 从备份数据库恢复区块, 共 %d 个区块，当前区块高度 %d, UUID %d\n", len(blocks), LocalChain.Head().Height, LocalChain.Head().UUID)
				} else {
					log.Printf("[ Backup] 从备份数据库恢复区块失败\n")
				}
			}
		} else {
			log.Printf("[ Backup] 无法从备份数据库读取区块: %v\n", err)
		}
		db.Close()
	} else {
		log.Printf("[ Backup] 无法打开备份数据库: %v\n", err)
	}

	discovery.Start(chaindb_config.DiscoveryInterval)

	validator.StartBlockDeliveryServer()
	validator.StartTryPickBlocks(LocalChain)

	blockchain.StartMining(LocalChain, &LocalTxPool, validator.Propose)

	monitor.Start(LocalChain, &LocalTxPool)
	StartUserService()

	go func() {
		for {
			time.Sleep(chaindb_config.BackupInterval)
			db, err := backup_db.New(chaindb_config.BackupDBFilePath)
			if err != nil {
				log.Printf("[ Backup] 打开备份数据失败%v\n", err)
				continue
			}
			defer db.Close()

			lastBlock := LocalChain.Head()
			uuidDBLast, err := db.GetLastBlockUUID()
			blocks := []blockchain.Block{}
			if err != nil {
				blocks = LocalChain.Snapshot()
			} else {
				for lastBlock.UUID != uuidDBLast {
					blocks = append(blocks, lastBlock)
					lastBlock, err = lastBlock.Previous()
					if err != nil {
						break
					}
				}
				slices.Reverse(blocks)
			}
			if len(blocks) > 0 {
				err = db.WriteBlocksToDB(blocks)
				if err != nil {
					log.Printf("[ Backup] 写入备份数据库失败: %v\n", err)
				} else {
					log.Printf("[ Backup] 成功备份区块到数据库, 共 %d 个区块\n", len(blocks))
				}
			}
		}
	}()

	select {}
}
