package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/me2/sms/rpc/internal/config"
	"github.com/me2/sms/rpc/internal/server"
	"github.com/me2/sms/rpc/internal/svc"
	"github.com/me2/sms/rpc/sms"
)

var configFile = flag.String("f", "etc/sms.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		sms.RegisterSmsServer(grpcServer, server.NewSmsServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting sms rpc server at %s...\n", c.ListenOn)
	s.Start()
}
