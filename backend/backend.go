package main

import (
	"flag"
	"fmt"
	"net/http"

	"model_mall_backend/backend/internal/config"
	"model_mall_backend/backend/internal/handler"
	aihandler "model_mall_backend/backend/internal/handler/ai"
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

	// 额外补充的AI接口（避免修改 goctl 生成的 routes.go）
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{ctx.JWTMiddleware},
			rest.Route{
				Method:  http.MethodPost,
				Path:    "/ai/anomaly/analyze",
				Handler: aihandler.AnalyzeAnomalyHandler(ctx),
			},
		),
		rest.WithPrefix("/api"),
	)

	// 静态文件服务已移除，如需要可以使用 Nginx 或其他静态文件服务器
	// 或者重新实现 StaticFileHandler

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
