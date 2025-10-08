package handler

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
	"modelmall/backend/internal/logic"
	"modelmall/backend/internal/svc"
	"modelmall/backend/internal/types"
)

// UploadImageHandler 处理图片上传
func UploadImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析multipart表单，最大32MB
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("解析表单失败: %w", err))
			return
		}

		// 获取上传的文件
		file, header, err := r.FormFile("image")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("获取上传文件失败: %w", err))
			return
		}
		defer file.Close()

		// 读取文件内容
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("读取文件内容失败: %w", err))
			return
		}

		// 检测图片尺寸
		_, err = file.Seek(0, 0)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("重置文件指针失败: %w", err))
			return
		}

		imgConfig, _, err := image.DecodeConfig(file)
		var width, height *int
		if err == nil {
			w := imgConfig.Width
			h := imgConfig.Height
			width = &w
			height = &h
		}

		// 创建上传目录
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("创建上传目录失败: %w", err))
			return
		}

		// 生成唯一文件名
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		filePath := filepath.Join(uploadDir, filename)

		// 保存文件
		outFile, err := os.Create(filePath)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("创建文件失败: %w", err))
			return
		}
		defer outFile.Close()

		if _, err := outFile.Write(fileBytes); err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("写入文件失败: %w", err))
			return
		}

		// 构造请求参数
		req := &logic.UploadImageRequest{
			Filename:   header.Filename,
			FilePath:   filePath,
			FileSize:   header.Size,
			MimeType:   header.Header.Get("Content-Type"),
			Width:      width,
			Height:     height,
			FileData:   fileBytes,
		}

		// 调用业务逻辑
		l := logic.NewUploadImageLogic(r.Context(), svcCtx)
		resp, err := l.UploadImage(req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// GetImageHandler 处理获取图片信息
func GetImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetImageReq
		
		// 从URL路径中获取id参数
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			// 尝试从路径参数获取
			idStr = r.PathValue("id")
		}
		
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, fmt.Errorf("无效的图片ID"))
			return
		}
		req.ImageID = id

		l := logic.NewGetImageLogic(r.Context(), svcCtx)
		resp, err := l.GetImage(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// ListImagesHandler 处理列出图片
func ListImagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListImagesReq
		
		// 解析查询参数
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewListImagesLogic(r.Context(), svcCtx)
		resp, err := l.ListImages(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}