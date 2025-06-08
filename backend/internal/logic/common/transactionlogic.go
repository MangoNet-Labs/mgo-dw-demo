package common

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
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
		db = db.Where("`from` = ?", authUser.MgoAddress)
	} else if req.Type == 2 {
		db = db.Where("`to` = ?", authUser.MgoAddress)
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
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 100
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
		return &types.TransactionListResp{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			List:     []types.MgoTransaction{},
		}, nil
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
		return &types.TransactionListResp{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			List:     []types.MgoTransaction{},
		}, nil
	}

	fmt.Println("signatures", signatures)

	return &types.TransactionListResp{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}, nil
}

type rawParsedTx struct {
	Result struct {
		Slot      uint64                  `json:"slot"`
		BlockTime *solana.UnixTimeSeconds `json:"blockTime"`
		Meta      struct {
			Fee         uint64   `json:"fee"`
			LogMessages []string `json:"logMessages"`
		} `json:"meta"`
		Transaction struct {
			Signatures []string `json:"signatures"`
			Message    struct {
				AccountKeys []struct {
					Pubkey string `json:"pubkey"`
				} `json:"accountKeys"`
				Instructions []struct {
					ProgramId string `json:"programId"`
					Parsed    struct {
						Type string                 `json:"type"`
						Info map[string]interface{} `json:"info"`
					} `json:"parsed"`
				} `json:"instructions"`
			} `json:"message"`
		} `json:"transaction"`
	} `json:"result"`
}
