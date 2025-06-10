package svc

import (
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/mangonet-labs/mgo-go-sdk/client"
	mgoConfig "github.com/mangonet-labs/mgo-go-sdk/config"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"time"
	"user/internal/config"
	"user/middleware"
)

type ServiceContext struct {
	Config         config.Config
	DB             *gorm.DB
	MgoCli         *client.Client
	SolCli         *rpc.Client
	AuthMiddleware rest.Middleware
	Client         *http.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DataSource), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	clientHttp := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &ServiceContext{
		Config:         c,
		DB:             db,
		MgoCli:         client.NewMgoClient(mgoConfig.RpcMgoTestnetEndpoint),
		SolCli:         rpc.New(rpc.DevNet_RPC),
		AuthMiddleware: middleware.NewJwtMiddleware(c.Auth.AccessSecret),
		Client:         clientHttp,
	}
}
