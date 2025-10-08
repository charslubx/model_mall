package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"model_mall_backend/backend/internal/svc"
)

// StaticFileHandler 静态文件处理器
func StaticFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取文件路径（去除 /uploads/ 前缀）
		filePath := r.URL.Path[len("/uploads/"):]
		
		// 构建完整的文件系统路径
		fullPath := filepath.Join(svcCtx.Config.Upload.StoragePath, filePath)
		
		// 检查文件是否存在
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		
		// 提供文件服务
		http.ServeFile(w, r, fullPath)
	}
}
