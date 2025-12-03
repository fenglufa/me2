package main

import (
	"flag"
	"fmt"

	"github.com/me2/scheduler/rpc/internal/config"
	schedcore "github.com/me2/scheduler/rpc/internal/scheduler"
	"github.com/me2/scheduler/rpc/internal/server"
	"github.com/me2/scheduler/rpc/internal/svc"
	"github.com/me2/scheduler/rpc/scheduler"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/scheduler.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// 启动调度器核心
	schedulerCore := schedcore.NewSchedulerCore(ctx)
	schedulerCore.Start()
	defer schedulerCore.Stop()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		scheduler.RegisterSchedulerServer(grpcServer, server.NewSchedulerServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	fmt.Printf("Scheduler core started with scan interval: %d seconds\n", c.Schedule.ScanInterval)
	s.Start()
}
