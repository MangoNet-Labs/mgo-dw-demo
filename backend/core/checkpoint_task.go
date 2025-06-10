package cron

import (
	"fmt"
	mgoModel "github.com/mangonet-labs/mgo-go-sdk/model"
	"log"
	"strings"
	"sync"
	"time"
	"user/internal/svc"
	"user/model"
	"user/third"
)

var taskRunning bool
var taskMutex sync.Mutex

func StartCheckpointSync(ctx *svc.ServiceContext) {

	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			taskMutex.Lock()
			if taskRunning {
				taskMutex.Unlock()
				continue // The previous task has not been completed yet
			}
			taskRunning = true
			taskMutex.Unlock()
			if err := syncCheckpoint(ctx); err != nil {
				log.Println("Checkpoint sync error:", err.Error())
			}
			taskMutex.Lock()
			taskRunning = false
			taskMutex.Unlock()
		}
	}()

}

func syncCheckpoint(ctx *svc.ServiceContext) error {

	sequenceNumber, err := third.GetLatestEpoch(ctx)
	if err != nil {
		sequenceNumber = "0"
	}
	fmt.Println("SequenceNumber:", sequenceNumber)
	resp, err := third.GetCheckpoints(sequenceNumber, ctx)
	if err != nil {
		return fmt.Errorf("failed to get checkpoints from sequence %s: %w", sequenceNumber, err)
	}

	if len(resp.Data) == 0 {
		return nil
	}
	checkpoints := make([]model.MgoCheckpoint, 0, len(resp.Data))
	for _, cp := range resp.Data {
		checkpoints = append(checkpoints, model.MgoCheckpoint{
			SequenceNumber:           cp.SequenceNumber,
			CheckpointDigest:         cp.Digest,
			Epoch:                    cp.Epoch,
			Transactions:             strings.Join(cp.Transactions, ","),
			PreviousCheckpointDigest: cp.PreviousDigest,
			TotalTransactionBlocks:   len(cp.Transactions),
			NetworkTotalTransactions: cp.NetworkTotalTransactions,
			TimestampMs:              cp.TimestampMs,
		})
	}

	if err := ctx.DB.CreateInBatches(checkpoints, len(checkpoints)).Error; err != nil {
		return fmt.Errorf("failed to save checkpoints to DB: %w", err)
	}

	return nil
}

func TransactionBlocksSync(ctx *svc.ServiceContext) {

	ticker := time.NewTicker(2 * time.Second)
	go func() {
		for range ticker.C {
			taskMutex.Lock()
			if taskRunning {
				taskMutex.Unlock()
				continue // The previous task has not been completed yet
			}
			taskRunning = true
			taskMutex.Unlock()
			if err := TransactionBlocks(ctx); err != nil {
				log.Println("Checkpoint sync error:", err.Error())
			}
			taskMutex.Lock()
			taskRunning = false
			taskMutex.Unlock()
		}
	}()

}

func TransactionBlocks(ctx *svc.ServiceContext) error {

	checkpoints, err := third.GetCheckpointList(ctx)
	if err != nil {
		return fmt.Errorf("failed to get checkpoint list: %w", err)
	}

	var updatedSeqs []string

	for _, checkpoint := range checkpoints {

		transactions, err := third.QueryTransactionBlocks(checkpoint.SequenceNumber, ctx)

		if err != nil {
			return fmt.Errorf("failed to query transaction blocks for checkpoint %d: %w", checkpoint.SequenceNumber, err)
		}
		if len(transactions) == 0 {
			continue
		}

		var mgoTxs []model.MgoTransaction

		for _, tx := range transactions {
			bc := tx.BalanceChanges

			// Only one-to-one MGO transfers are processed
			if len(bc) < 2 || !allCoinTypeMatch(bc, "0x2::mgo::MGO") {
				continue
			}

			var from, to, amount, fromAmount string

			if len(bc) == 2 {

				for _, change := range bc {
					if strings.HasPrefix(change.Amount, "-") {
						from = change.Owner.AddressOwner
						fromAmount = change.Amount
					} else {
						to = change.Owner.AddressOwner
						amount = change.Amount
					}
				}

				mgoTxs = append(mgoTxs, model.MgoTransaction{
					Digest:      tx.Digest,
					From:        from,
					To:          to,
					Amount:      amount,
					FromAmount:  fromAmount,
					Checkpoint:  tx.Checkpoint,
					CoinType:    bc[0].CoinType,
					GasOwner:    tx.Transaction.Data.GasData.Owner,
					GasPrice:    tx.Transaction.Data.GasData.Price,
					GasBudget:   tx.Transaction.Data.GasData.Budget,
					TimestampMs: tx.TimestampMs,
				})
			} else {

				for _, change := range bc {
					if strings.HasPrefix(change.Amount, "-") {
						if change.Owner.AddressOwner == tx.Transaction.Data.GasData.Owner {
							continue
						}
						from = change.Owner.AddressOwner
						fromAmount = change.Amount
					} else {
						to = change.Owner.AddressOwner
						amount = change.Amount
					}
				}

				mgoTxs = append(mgoTxs, model.MgoTransaction{
					Digest:      tx.Digest,
					From:        from,
					To:          to,
					Amount:      amount,
					FromAmount:  fromAmount,
					Checkpoint:  tx.Checkpoint,
					CoinType:    bc[0].CoinType,
					GasOwner:    tx.Transaction.Data.GasData.Owner,
					GasPrice:    tx.Transaction.Data.GasData.Price,
					GasBudget:   tx.Transaction.Data.GasData.Budget,
					TimestampMs: tx.TimestampMs,
				})
			}

		}

		if len(mgoTxs) > 0 {
			if err := ctx.DB.CreateInBatches(mgoTxs, len(mgoTxs)).Error; err != nil {
				return fmt.Errorf("failed to save mgo transactions: %w", err)
			}
		}
		updatedSeqs = append(updatedSeqs, checkpoint.SequenceNumber)

	}

	if len(updatedSeqs) > 0 {
		err := ctx.DB.Model(&model.MgoCheckpoint{}).
			Where("sequence_number IN ?", updatedSeqs).
			Update("status", 1).Error
		if err != nil {
			return fmt.Errorf("batch update checkpoint status failed: %w", err)
		}
	}
	return nil

}

func allCoinTypeMatch(balanceChanges []mgoModel.BalanceChanges, target string) bool {
	for _, change := range balanceChanges {
		if change.CoinType != target {
			return false
		}
	}
	return true
}
