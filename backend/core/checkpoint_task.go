package cron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	mgoModel "github.com/mangonet-labs/mgo-go-sdk/model"
	"log"
	"net/http"
	"strconv"
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

func TransactionSignatureSync(ctx *svc.ServiceContext) {

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
			if err := TransactionSignatureBlocks(ctx); err != nil {
				log.Println("Checkpoint sync error:", err.Error())
			}
			taskMutex.Lock()
			taskRunning = false
			taskMutex.Unlock()
		}
	}()

}

func TransactionSignatureBlocks(ctx *svc.ServiceContext) error {

	var checkpoints []model.SolTransaction
	err := ctx.DB.Select("id,digest,status,gas_owner").Where("status = 0 ").Limit(10).Find(&checkpoints).Error
	if err != nil {
		return err
	}
	url := ctx.Config.HeliusRpc
	ctxBack := context.Background()
	for _, signature := range checkpoints {

		reqBody := jsonrpc.RPCRequest{
			JSONRPC: "2.0",
			ID:      "1",
			Method:  "getTransaction",
			Params: []interface{}{
				signature.Digest,
				"jsonParsed",
			},
		}
		payload, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println("json marshal error:", err)
			continue
		}
		reqResult, err := http.NewRequestWithContext(ctxBack, "POST", url, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("http request error:", err)
			continue
		}
		reqResult.Header.Set("Content-Type", "application/json")
		transactionResult, err := ctx.Client.Do(reqResult)
		if err != nil {
			fmt.Println("http request error:", err)
			continue
		}
		defer transactionResult.Body.Close()
		var rpcResp Transaction
		if err := json.NewDecoder(transactionResult.Body).Decode(&rpcResp); err != nil {
			fmt.Println("json unmarshal error:", err)
			continue
		}

		if len(rpcResp.Result.Meta.InnerInstructions) > 0 {
			for _, innerInstruction := range rpcResp.Result.Meta.InnerInstructions {
				for _, instruction := range innerInstruction.Instructions {
					if instruction.Parsed.Type == "mintTo" {
						info := instruction.Parsed.Info
						tx := model.SolTransaction{
							From:     info.Mint,
							To:       signature.GasOwner,
							Amount:   info.Amount,
							CoinType: ctx.Config.SplToken,
							GasPrice: strconv.Itoa(rpcResp.Result.Meta.Fee),
							Status:   1,
						}
						ctx.DB.Where("id = ?", signature.ID).Updates(&tx)
					}
				}
			}
		}

		for _, instr := range rpcResp.Result.Transaction.Message.Instructions {

			if instr.Parsed.Type == "transfer" {
				info := instr.Parsed.Info
				To := info.Destination
				if len(rpcResp.Result.Meta.PostTokenBalances) == 1 {
					To = rpcResp.Result.Meta.PostTokenBalances[0].Owner
				} else if len(rpcResp.Result.Meta.PostTokenBalances) == 2 {
					To = rpcResp.Result.Meta.PostTokenBalances[1].Owner
				}
				tx := model.SolTransaction{
					From:     info.Authority,
					To:       To,
					Amount:   info.Amount,
					CoinType: ctx.Config.SplToken,
					GasPrice: strconv.Itoa(rpcResp.Result.Meta.Fee),
					Status:   1,
				}
				ctx.DB.Where("id = ?", signature.ID).Updates(&tx)
			} else if instr.Parsed.Type == "transferChecked" {

				info := instr.Parsed.Info
				To := info.Destination
				if len(rpcResp.Result.Meta.PostTokenBalances) == 1 {
					To = rpcResp.Result.Meta.PostTokenBalances[0].Owner
				} else if len(rpcResp.Result.Meta.PostTokenBalances) == 2 {
					To = rpcResp.Result.Meta.PostTokenBalances[1].Owner
				}
				tx := model.SolTransaction{
					From:     info.MultisigAuthority,
					To:       To,
					Amount:   info.TokenAmount.Amount,
					CoinType: ctx.Config.SplToken,
					GasPrice: strconv.Itoa(rpcResp.Result.Meta.Fee),
					Status:   1,
				}
				ctx.DB.Where("id = ?", signature.ID).Updates(&tx)
			}
		}
		ctx.DB.Where("id = ?", signature.ID).Update("status", "1")

	}

	return nil

}

type Transaction struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BlockTime int `json:"blockTime"`
		Meta      struct {
			ComputeUnitsConsumed int         `json:"computeUnitsConsumed"`
			CostUnits            int         `json:"costUnits"`
			Err                  interface{} `json:"err"`
			Fee                  int         `json:"fee"`
			InnerInstructions    []struct {
				Index        int `json:"index"`
				Instructions []struct {
					Parsed struct {
						Info struct {
							Account       string `json:"account"`
							Amount        string `json:"amount"`
							Mint          string `json:"mint"`
							MintAuthority string `json:"mintAuthority"`
						} `json:"info"`
						Type string `json:"type"`
					} `json:"parsed"`
					Program     string `json:"program"`
					ProgramId   string `json:"programId"`
					StackHeight int    `json:"stackHeight"`
				} `json:"instructions"`
			} `json:"innerInstructions"`
			LogMessages       []string `json:"logMessages"`
			PostBalances      []int64  `json:"postBalances"`
			PostTokenBalances []struct {
				AccountIndex  int    `json:"accountIndex"`
				Mint          string `json:"mint"`
				Owner         string `json:"owner"`
				ProgramId     string `json:"programId"`
				UiTokenAmount struct {
					Amount         string  `json:"amount"`
					Decimals       int     `json:"decimals"`
					UiAmount       float64 `json:"uiAmount"`
					UiAmountString string  `json:"uiAmountString"`
				} `json:"uiTokenAmount"`
			} `json:"postTokenBalances"`
			PreBalances      []int64 `json:"preBalances"`
			PreTokenBalances []struct {
				AccountIndex  int    `json:"accountIndex"`
				Mint          string `json:"mint"`
				Owner         string `json:"owner"`
				ProgramId     string `json:"programId"`
				UiTokenAmount struct {
					Amount         string  `json:"amount"`
					Decimals       int     `json:"decimals"`
					UiAmount       float64 `json:"uiAmount"`
					UiAmountString string  `json:"uiAmountString"`
				} `json:"uiTokenAmount"`
			} `json:"preTokenBalances"`
			Rewards []interface{} `json:"rewards"`
			Status  struct {
				Ok interface{} `json:"Ok"`
			} `json:"status"`
		} `json:"meta"`
		Slot        int `json:"slot"`
		Transaction struct {
			Message struct {
				AccountKeys []struct {
					Pubkey   string `json:"pubkey"`
					Signer   bool   `json:"signer"`
					Source   string `json:"source"`
					Writable bool   `json:"writable"`
				} `json:"accountKeys"`
				Instructions []struct {
					Parsed struct {
						Info struct {
							Amount            string   `json:"amount"`
							Authority         string   `json:"authority"`
							Destination       string   `json:"destination"`
							Source            string   `json:"source"`
							Mint              string   `json:"mint"`
							MultisigAuthority string   `json:"multisigAuthority"`
							Signers           []string `json:"signers"`
							TokenAmount       struct {
								Amount         string  `json:"amount"`
								Decimals       int     `json:"decimals"`
								UiAmount       float64 `json:"uiAmount"`
								UiAmountString string  `json:"uiAmountString"`
							} `json:"tokenAmount"`
						} `json:"info"`
						Type string `json:"type"`
					} `json:"parsed"`
					Program     string `json:"program"`
					ProgramId   string `json:"programId"`
					StackHeight int    `json:"stackHeight"`
				} `json:"instructions"`
				RecentBlockhash string `json:"recentBlockhash"`
			} `json:"message"`
			Signatures []string `json:"signatures"`
		} `json:"transaction"`
	} `json:"result"`
	Id string `json:"id"`
}
