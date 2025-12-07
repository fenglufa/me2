package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/me2/gateway/api/internal/config"
	"github.com/me2/gateway/api/internal/handler"
	dialogueHandler "github.com/me2/gateway/api/internal/handler/dialogue"
	"github.com/me2/gateway/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 手动注册 WebSocket 路由
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/dialogue/stream",
		Handler: dialogueHandler.StreamHandlerWithAuth(ctx),
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
