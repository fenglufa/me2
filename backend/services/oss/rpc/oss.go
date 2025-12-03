package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/me2/oss/rpc/internal/config"
	"github.com/me2/oss/rpc/internal/server"
	"github.com/me2/oss/rpc/internal/svc"
	"github.com/me2/oss/rpc/oss"
)

var configFile = flag.String("f", "etc/oss.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		oss.RegisterOssServer(grpcServer, server.NewOssServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting oss rpc server at %s...\n", c.ListenOn)
	s.Start()
}
