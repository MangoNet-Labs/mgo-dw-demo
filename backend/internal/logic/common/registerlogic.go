package common

import (
	"context"
	"errors"
	"log"
	"user/common"
	"user/internal/svc"
	"user/internal/types"
	"user/model"

	"github.com/gagliardetto/solana-go"
	"github.com/mangonet-labs/mgo-go-sdk/account/keypair"
	"github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.UserResponse, err error) {
	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}
	var exists model.Address
	tx := l.svcCtx.DB.Where("username = ?", req.Username).First(&exists)
	if tx.Error == nil {
		return nil, errors.New("username already exists")
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	// Mango SDK Generate a Key Pair
	kp, err := keypair.NewKeypair(config.Ed25519Flag)
	if err != nil {
		log.Fatalf("failed to generate keypair: %v", err)
	}
	//Sol Go SDK
	wallet := solana.NewWallet()
	user := model.Address{
		Username:         req.Username,
		Password:         string(hashed),
		MgoAddress:       kp.MgoAddress(),
		MgoPrivateKey:    kp.PrivateKeyHex(),
		SolanaAddress:    wallet.PublicKey().String(),
		SolanaPrivateKey: wallet.PrivateKey.String(),
	}
	if err := l.svcCtx.DB.Create(&user).Error; err != nil {
		return nil, err
	}
	token, err := common.GenerateToken(user, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		return nil, err
	}
	return &types.UserResponse{
		Id:            user.ID,
		Username:      user.Username,
		MgoAddress:    user.MgoAddress,
		SolanaAddress: user.SolanaAddress,
		Token:         token,
	}, nil
}
