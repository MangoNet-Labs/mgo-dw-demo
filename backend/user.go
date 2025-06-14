package main

import (
	"flag"
	"fmt"
	cron "user/core"
	"user/internal/config"
	"user/internal/handler"
	"user/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {

	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	server := rest.MustNewServer(c.RestConf, rest.WithCors("*"))
	defer server.Stop()
	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	//Storage Checkpoints
	cron.StartCheckpointSync(ctx)
	//Query checkpoint transactions
	cron.TransactionBlocksSync(ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
