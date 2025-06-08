package common

import (
	"context"
	"errors"
	"fmt"
	bin "github.com/gagliardetto/binary"
	"math/big"
	"strconv"
	"user/internal/svc"
	"user/internal/types"
	"user/middleware"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/model/response"

	tokenSol "github.com/gagliardetto/solana-go/programs/token"
	"github.com/zeromicro/go-zero/core/logx"
)

type BalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BalanceLogic {
	return &BalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BalanceLogic) Balance(req *types.GetBalanceByCoinReq, authUser *middleware.AuthUser) (res *response.CoinBalanceResponse, err error) {

	if req.ChainName == "mgo" {

		CoinBalanceResponse, err := l.svcCtx.MgoCli.MgoXGetBalance(l.ctx, request.MgoXGetBalanceRequest{
			Owner:    authUser.MgoAddress,
			CoinType: "0x2::mgo::MGO",
		})
		if err != nil {
			return nil, err
		}
		CoinBalanceResponse.TotalBalance = ParseBalance(CoinBalanceResponse.TotalBalance)
		return &CoinBalanceResponse, nil

	} else if req.ChainName == "sol" {

		token := l.svcCtx.Config.SplToken
		mintPub, err := solana.PublicKeyFromBase58(token)
		publicKey := solana.MustPublicKeyFromBase58(authUser.SolAddress)
		out, err := l.svcCtx.SolCli.GetTokenAccountsByOwner(l.ctx, publicKey, &rpc.GetTokenAccountsConfig{
			Mint: &mintPub,
		},
			&rpc.GetTokenAccountsOpts{
				Encoding: solana.EncodingBase64Zstd,
			})
		if err != nil {
			return nil, err
		}

		totalRaw := big.NewInt(0)
		for _, rawAccount := range out.Value {
			var tokAcc tokenSol.Account
			data := rawAccount.Account.Data.GetBinary()
			dec := bin.NewBinDecoder(data)
			err := dec.Decode(&tokAcc)
			if err != nil {
				fmt.Println("decode error:", err)
				continue
			}
			amt := big.NewInt(0).SetUint64(tokAcc.Amount)
			totalRaw.Add(totalRaw, amt)
		}

		return &response.CoinBalanceResponse{
			CoinType:        l.svcCtx.Config.SplToken,
			TotalBalance:    ParseBalance(strconv.FormatUint(totalRaw.Uint64(), 10)),
			CoinObjectCount: int(out.Context.Slot),
		}, nil

	} else {
		return nil, errors.New("parameter error")
	}
}

func ParseBalance(balanceStr string) string {
	balance := new(big.Int)
	balance.SetString(balanceStr, 10)
	if balance.Cmp(big.NewInt(0)) > 0 {
		balanceFloat := new(big.Float).SetInt(balance)
		divider := new(big.Float).SetFloat64(1e9)
		result := new(big.Float).Quo(balanceFloat, divider)
		return result.Text('f', 9)
	}
	return "0"
}
