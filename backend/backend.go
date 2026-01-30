package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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

	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// 确保上传目录存在（静态文件由Nginx提供服务）
	ensureUploadDirectory(c)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	fmt.Println("Note: Static files should be served by Nginx at /uploads/*")
	server.Start()
}

// ensureUploadDirectory 确保上传目录存在
func ensureUploadDirectory(c config.Config) {
	uploadDir := c.Upload.StoragePath
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// 转换为绝对路径
	absUploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		fmt.Printf("Warning: Failed to get absolute path for upload directory: %v\n", err)
		absUploadDir = uploadDir
	}

	// 创建上传目录及子目录
	dirs := []string{
		absUploadDir,
		filepath.Join(absUploadDir, "avatar"),
		filepath.Join(absUploadDir, "product"),
		filepath.Join(absUploadDir, "temp"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Warning: Failed to create directory %s: %v\n", dir, err)
		}
	}

	fmt.Printf("Upload directory initialized: %s\n", absUploadDir)
}
