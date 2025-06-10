package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"net/http"
	"strconv"
	"user/middleware"

	"user/internal/svc"
	"user/internal/types"
	"user/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransactionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionLogic {
	return &TransactionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransactionLogic) Transaction(req *types.TransactionListReq, authUser *middleware.AuthUser) (resp *types.TransactionListResp, err error) {

	var list []types.MgoTransaction
	var total int64
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	db := l.svcCtx.DB.Model(&model.MgoTransaction{})
	if req.Type == 1 {
		db = db.Where("`to` = ?", authUser.MgoAddress)
	} else if req.Type == 2 {
		db = db.Where("`from` = ?", authUser.MgoAddress)
	} else {
		db = db.Where("`from` = ? OR `to` = ?", authUser.MgoAddress, authUser.MgoAddress)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := db.Order("timestamp_ms DESC").Limit(pageSize).
		Offset(offset).Find(&list).Error; err != nil {
		return nil, err
	}

	return &types.TransactionListResp{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}, nil
}

func (l *TransactionLogic) TransactionSol(req *types.TransactionListReq, authUser *middleware.AuthUser) (resp *types.TransactionListResp, err error) {

	var list []types.MgoTransaction
	var total int64
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 6 {
		pageSize = 6
	}
	userPubKey := solana.MustPublicKeyFromBase58(authUser.SolAddress)
	tokenMint := solana.MustPublicKeyFromBase58(l.svcCtx.Config.SplToken)
	tokenAccountsResult, err := l.svcCtx.SolCli.GetTokenAccountsByOwner(
		l.ctx,
		userPubKey,
		&rpc.GetTokenAccountsConfig{Mint: &tokenMint},
		&rpc.GetTokenAccountsOpts{},
	)
	if err != nil {
		return nil, err
	}
	if len(tokenAccountsResult.Value) == 0 {
		return nil, nil
	}
	tokenAccountPubKey := tokenAccountsResult.Value[0].Pubkey
	signatures, err := l.svcCtx.SolCli.GetSignaturesForAddressWithOpts(l.ctx, tokenAccountPubKey, &rpc.GetSignaturesForAddressOpts{
		Limit:      &pageSize,
		Commitment: rpc.CommitmentFinalized,
	})
	if err != nil {
		l.Logger.Errorf("GetSignaturesForAddress error: %v", err)
		return nil, err
	}
	totalSigs := len(signatures)
	if totalSigs == 0 {
		return nil, nil
	}

	url := l.svcCtx.Config.HeliusRpc
	for _, signature := range signatures {

		reqBody := jsonrpc.RPCRequest{
			JSONRPC: "2.0",
			ID:      "1",
			Method:  "getTransaction",
			Params: []interface{}{
				signature.Signature,
				"jsonParsed",
			},
		}
		payload, err := json.Marshal(reqBody)
		if err != nil {
			fmt.Println("json marshal error:", err)
			continue
		}
		reqResult, err := http.NewRequestWithContext(l.ctx, "POST", url, bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("http request error:", err)
			continue
		}
		reqResult.Header.Set("Content-Type", "application/json")
		transactionResult, err := l.svcCtx.Client.Do(reqResult)
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
		timestampMs := fmt.Sprintf("%d", rpcResp.Result.BlockTime*1000)
		for _, instr := range rpcResp.Result.Transaction.Message.Instructions {
			if instr.Parsed.Type != "transfer" {
				continue
			}
			info := instr.Parsed.Info
			tx := types.MgoTransaction{
				ID:          int64(rpcResp.Result.Slot),
				Digest:      signature.Signature.String(),
				From:        info.Authority,
				To:          info.Destination,
				Amount:      info.Amount,
				TimestampMs: timestampMs,
				Checkpoint:  strconv.Itoa(rpcResp.Result.Slot),
				CoinType:    l.svcCtx.Config.SplToken,
				GasOwner:    info.Source,
				GasPrice:    strconv.Itoa(rpcResp.Result.Meta.Fee),
			}

			if req.Type == 1 {
				if tx.From != authUser.SolAddress {
					list = append(list, tx)
				}
			} else if req.Type == 2 {
				if tx.From == authUser.SolAddress {
					list = append(list, tx)
				}
			} else {
				list = append(list, tx)
			}

		}
	}
	if len(list) == 0 {
		list = []types.MgoTransaction{}
	}
	return &types.TransactionListResp{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}, nil
}

type Transaction struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		BlockTime int `json:"blockTime"`
		Meta      struct {
			ComputeUnitsConsumed int           `json:"computeUnitsConsumed"`
			CostUnits            int           `json:"costUnits"`
			Err                  interface{}   `json:"err"`
			Fee                  int           `json:"fee"`
			InnerInstructions    []interface{} `json:"innerInstructions"`
			LogMessages          []string      `json:"logMessages"`
			PostBalances         []int64       `json:"postBalances"`
			PostTokenBalances    []struct {
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
							Amount      string `json:"amount"`
							Authority   string `json:"authority"`
							Destination string `json:"destination"`
							Source      string `json:"source"`
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
