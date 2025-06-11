package common

import (
	"context"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
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

func (l *TransactionLogic) TransactionSol(req *types.TransactionListReq, authUser *middleware.AuthUser) (resp *types.TransactionSolResp, err error) {

	var list []types.SolTransaction
	var total int64
	page := req.Page
	pageSize := req.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 10 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	db := l.svcCtx.DB.Model(&model.SolTransaction{})
	switch req.Type {
	case 1:
		db.Where("`to` = ?", authUser.SolAddress)
	case 2:
		db.Where("`from` = ?", authUser.SolAddress)
	default:
		db.Where("`from` = ? OR `to` = ?", authUser.SolAddress, authUser.SolAddress)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if err := db.Order("checkpoint DESC").Limit(pageSize).
		Offset(offset).Find(&list).Error; err != nil {
		return nil, err
	}
	_, err = l.fetchLatestSolTransactions(authUser)
	if err != nil {
		fmt.Println("err:", err)
	}
	return &types.TransactionSolResp{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		List:     list,
	}, nil

}

func (l *TransactionLogic) fetchLatestSolTransactions(
	authUser *middleware.AuthUser,
) ([]types.SolTransaction, error) {

	var checkpoint types.SolTransaction
	pageSize := 50
	reqSign := rpc.GetSignaturesForAddressOpts{
		Limit:      &pageSize,
		Commitment: rpc.CommitmentFinalized,
	}
	l.svcCtx.DB.Where("gas_owner = ? ", authUser.SolAddress).
		Order("checkpoint DESC").Limit(1).Find(&checkpoint)
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
	signatures, err := l.svcCtx.SolCli.GetSignaturesForAddressWithOpts(l.ctx, tokenAccountPubKey, &reqSign)
	if err != nil {
		l.Logger.Errorf("GetSignaturesForAddress error: %v", err)
		return nil, err
	}
	totalSigs := len(signatures)
	if totalSigs == 0 {
		return nil, nil
	}
	for _, signature := range signatures {
		if signature.Slot > checkpoint.Checkpoint {
			l.svcCtx.DB.Create(&types.SolTransaction{
				Digest:      signature.Signature.String(),
				TimestampMs: strconv.FormatInt(signature.BlockTime.Time().UnixMilli(), 10),
				Checkpoint:  signature.Slot,
				GasOwner:    authUser.SolAddress,
			})
		}
	}
	return nil, nil
}
