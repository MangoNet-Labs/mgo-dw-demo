package common

import (
	"context"
	"errors"
	"user/common"
	"user/internal/svc"
	"user/internal/types"
	"user/model"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.UserResponse, err error) {
	var user model.Address
	tx := l.svcCtx.DB.Where("username = ?", req.Username).First(&user)
	if tx.Error != nil {
		return nil, errors.New("user does not exist")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("wrong password")
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
