package main

import (
	"flag"
	"fmt"
	"net/http"

	"model_mall_backend/backend/internal/config"
	"model_mall_backend/backend/internal/handler"
	"model_mall_backend/backend/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/backend-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 注册静态文件服务（用于访问上传的图片）
	// 使用通配符路径来支持嵌套目录
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/uploads/",
		Handler: handler.StaticFileHandler(ctx),
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
