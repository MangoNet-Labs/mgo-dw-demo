package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	mgoModel "github.com/mangonet-labs/mgo-go-sdk/model"
	"github.com/mangonet-labs/mgo-go-sdk/model/request"
	"github.com/mangonet-labs/mgo-go-sdk/transaction"
	"github.com/shopspring/decimal"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"user/internal/svc"
	"user/internal/types"
	"user/middleware"
	"user/model"
)

type WithdrawalLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWithdrawalLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawalLogic {
	return &WithdrawalLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WithdrawalLogic) Withdrawal(req *types.WithdrawalReq, authUser *middleware.AuthUser, TotalBalance string) (resp *types.WithdrawalResp, err error) {

	if len(req.ToAddress) != 66 {
		return nil, errors.New("receiving address error")
	}
	var user model.Address
	txUser := l.svcCtx.DB.Where("id = ?", authUser.UserID).First(&user)
	if txUser.Error != nil {
		return nil, errors.New("user does not exist")
	}
	Balance, err := strconv.ParseFloat(TotalBalance, 10)
	if err != nil {
		return nil, err
	}
	if Balance < req.Amount {
		return nil, errors.New("insufficient cash withdrawal amount")
	}
	cli := l.svcCtx.MgoCli
	var ctx = context.Background()
	walletKey, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, user.MgoPrivateKey)
	if err != nil {
		return nil, err
	}
	sysKey, err := keypair.NewKeypairWithPrivateKey(config.Ed25519Flag, l.svcCtx.Config.SysMgoPrivateKey)
	if err != nil {
		return nil, err
	}
	gasCoinObj, err := cli.MgoGetObject(ctx, request.MgoGetObjectRequest{ObjectId: l.svcCtx.Config.SysGasObject})
	if err != nil {
		return nil, err
	}
	gasCoin, err := transaction.NewMgoObjectRef(
		mgoModel.MgoAddress(gasCoinObj.Data.ObjectId),
		gasCoinObj.Data.Version,
		mgoModel.ObjectDigest(gasCoinObj.Data.Digest),
	)
	if err != nil {
		return nil, err
	}
	tx := transaction.NewTransaction()
	tx.SetMgoClient(cli).
		SetSigner(walletKey).
		SetSponsoredSigner(sysKey).
		SetSender(mgoModel.MgoAddress(walletKey.MgoAddress())).
		SetGasPrice(1000).
		SetGasBudget(50000000).
		SetGasPayment([]transaction.MgoObjectRef{*gasCoin}).
		SetGasOwner(mgoModel.MgoAddress(sysKey.MgoAddress()))

	amountDecimal := decimal.NewFromFloat(req.Amount)
	amount := amountDecimal.Mul(decimal.New(1, 9)).BigInt().Uint64()
	mergeCoin, err := l.GetEnoughMgo(ctx, walletKey.MgoAddress(), amountDecimal.Mul(decimal.New(1, 9)), tx)
	if err != nil {
		return nil, err
	}
	splitCoin := tx.SplitCoins(mergeCoin, []transaction.Argument{
		tx.Pure(amount),
	})
	tx.TransferObjects([]transaction.Argument{splitCoin}, tx.Pure(req.ToAddress))

	MgoTransactionBlockResponse, err := tx.Execute(
		ctx,
		request.MgoTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
			ShowEvents:   true,
		},
		"WaitForLocalExecution",
	)
	if err != nil {
		return nil, err
	}
	return &types.WithdrawalResp{
		Hash: MgoTransactionBlockResponse.Digest,
	}, nil
}

func (l *WithdrawalLogic) GetEnoughMgo(ctx context.Context, owner string, amount decimal.Decimal, tx *transaction.Transaction) (mergeCoin transaction.Argument, err error) {

	fmt.Println(amount)
	var sumBanlance decimal.Decimal
	for {
		paginatedCoins, err := l.svcCtx.MgoCli.MgoXGetCoins(ctx, request.MgoXGetCoinsRequest{
			Owner:    owner,
			CoinType: "0x2::mgo::MGO",
			Limit:    50,
		})
		if err != nil {
			return mergeCoin, err
		}
		if len(paginatedCoins.Data) == 0 {
			return mergeCoin, errors.New("no enough mgo")
		}
		for _, coin := range paginatedCoins.Data {

			if coin.Balance == "0" {
				continue
			}
			coinBanlance, err := decimal.NewFromString(coin.Balance)
			if err != nil {
				return mergeCoin, err
			}
			sumBanlance = sumBanlance.Add(coinBanlance)
			gasCoin, err := transaction.NewMgoObjectRef(
				mgoModel.MgoAddress(coin.CoinObjectId),
				coin.Version,
				mgoModel.ObjectDigest(coin.Digest),
			)
			if mergeCoin.Input == nil {
				mergeCoin = tx.Object(
					transaction.CallArg{Object: &transaction.ObjectArg{ImmOrOwnedObject: gasCoin}},
				)
			} else {
				mergeCoin = tx.MergeCoins(mergeCoin,
					[]transaction.Argument{
						tx.Object(transaction.CallArg{Object: &transaction.ObjectArg{ImmOrOwnedObject: gasCoin}}),
					})
			}
			if sumBanlance.Cmp(amount) >= 0 {
				return mergeCoin, err
			}
		}

	}

}

func (l *WithdrawalLogic) WithdrawalSol(req *types.WithdrawalReq, authUser *middleware.AuthUser, TotalBalance string) (resp *types.WithdrawalResp, err error) {

	if len(req.ToAddress) < 42 {
		return nil, errors.New("invalid receiving address")
	}

	balance, err := strconv.ParseFloat(TotalBalance, 10)
	if err != nil {
		return nil, fmt.Errorf("parse balance error: %w", err)
	}
	if balance < req.Amount {
		return nil, errors.New("insufficient withdrawal amount")
	}

	amount := uint64(req.Amount * 1e9) // Convert to lamports (assuming token has 9 decimals)

	// Fetch user from DB
	var user model.Address
	if err := l.svcCtx.DB.Where("id = ?", authUser.UserID).First(&user).Error; err != nil {
		return nil, errors.New("user does not exist")
	}

	cli := l.svcCtx.SolCli
	splToken := l.svcCtx.Config.SplToken

	// Parse system private key and user wallet private key
	sysPriv := solana.MustPrivateKeyFromBase58(l.svcCtx.Config.SysSolPrivateKey)
	sysPub := sysPriv.PublicKey()

	wallet, err := solana.PrivateKeyFromBase58(user.SolanaPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid user private key: %w", err)
	}

	// Parse SPL token mint and recipient address
	tokenMint := solana.MustPublicKeyFromBase58(splToken)
	toPubkey, err := solana.PublicKeyFromBase58(req.ToAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient address: %w", err)
	}

	// Find sender's and recipient's associated token accounts (ATA)
	senderATA, _, err := solana.FindAssociatedTokenAddress(wallet.PublicKey(), tokenMint)
	if err != nil {
		return nil, fmt.Errorf("find sender ATA error: %w", err)
	}
	receiverATA, _, err := solana.FindAssociatedTokenAddress(toPubkey, tokenMint)
	if err != nil {
		return nil, fmt.Errorf("find receiver ATA error: %w", err)
	}

	// Check if recipient's ATA exists
	receiverExists, err := l.tokenAccountExists(context.Background(), receiverATA)
	if err != nil {
		return nil, fmt.Errorf("check receiver ATA error: %w", err)
	}

	// Get recent blockhash for transaction
	respBlock, err := cli.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return nil, fmt.Errorf("get latest blockhash error: %w", err)
	}

	var instructions []solana.Instruction

	// If receiver ATA does not exist, create it, payer is system account
	if !receiverExists {
		createATAInstruction := solana.NewInstruction(
			solana.SPLAssociatedTokenAccountProgramID,
			solana.AccountMetaSlice{
				// Payer: system account pays rent and signs
				{PublicKey: sysPub, IsSigner: true, IsWritable: true},
				// ATA to create
				{PublicKey: receiverATA, IsSigner: false, IsWritable: true},
				// Owner of new ATA
				{PublicKey: toPubkey, IsSigner: false, IsWritable: false},
				// Mint
				{PublicKey: tokenMint, IsSigner: false, IsWritable: false},
				// System program
				{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
				// Token program
				{PublicKey: solana.TokenProgramID, IsSigner: false, IsWritable: false},
				// Rent sysvar
				{PublicKey: solana.SysVarRentPubkey, IsSigner: false, IsWritable: false},
			},
			nil,
		)
		instructions = append(instructions, createATAInstruction)
	}

	// Build SPL token transfer instruction
	transferInstruction := token.NewTransferInstruction(
		amount,
		senderATA,
		receiverATA,
		wallet.PublicKey(),
		[]solana.PublicKey{},
	).Build()
	instructions = append(instructions, transferInstruction)

	// Build transaction with system account as payer (to pay for ATA creation)
	tx, err := solana.NewTransaction(
		instructions,
		respBlock.Value.Blockhash,
		solana.TransactionPayer(sysPub),
	)
	if err != nil {
		return nil, fmt.Errorf("create transaction error: %w", err)
	}

	// Sign transaction with both system account and user wallet
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		switch {
		case key.Equals(sysPub):
			return &sysPriv
		case key.Equals(wallet.PublicKey()):
			return &wallet
		default:
			return nil
		}
	})
	if err != nil {
		return nil, fmt.Errorf("sign transaction error: %w", err)
	}

	// Send transaction
	sig, err := cli.SendTransaction(context.Background(), tx)
	if err != nil {
		return nil, fmt.Errorf("send transaction error: %w", err)
	}

	return &types.WithdrawalResp{
		Hash: sig.String(),
	}, nil

}

func (l *WithdrawalLogic) tokenAccountExists(ctx context.Context, account solana.PublicKey) (bool, error) {
	_, err := l.svcCtx.SolCli.GetAccountInfo(ctx, account)
	if err != nil {
		return false, err
	}
	return true, nil
}
