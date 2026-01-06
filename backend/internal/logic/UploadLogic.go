package logic

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

const propConfigMaxSize = int64(10 * 1024 * 1024) // 10MB

type UploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadLogic {
	return &UploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadFile 上传通用文件（大小限制由 Upload.MaxSize 控制）
func (l *UploadLogic) UploadFile(r *http.Request) (*types.UploadFileResp, error) {
	userID := l.getUserIDFromContext()
	if userID == 0 {
		return nil, fmt.Errorf("未授权")
	}

	maxSize := l.svcCtx.Config.Upload.MaxSize
	if maxSize <= 0 {
		maxSize = propConfigMaxSize
	}

	if err := r.ParseMultipartForm(maxSize); err != nil {
		return nil, fmt.Errorf("解析表单失败: %w", err)
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("获取上传文件失败: %w", err)
	}
	defer file.Close()

	if fileHeader.Size > maxSize {
		return nil, fmt.Errorf("文件大小超过限制")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	dateDir := time.Now().Format("2006/01/02")
	storageDir := filepath.Join(l.svcCtx.Config.Upload.StoragePath, "files", dateDir)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	dstPath := filepath.Join(storageDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		mimeType = http.DetectContentType(buf[:n])
		_, _ = file.Seek(0, 0)
	}

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s/%s", l.svcCtx.Config.Upload.BaseURL, "files", dateDir, filename)
	return &types.UploadFileResp{
		Url:      url,
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		MimeType: mimeType,
	}, nil
}

// UploadPropConfig 上传属性配置文件（固定 <= 10MB）
// 允许的扩展名：.json / .yaml / .yml
func (l *UploadLogic) UploadPropConfig(r *http.Request) (*types.UploadPropConfigResp, error) {
	userID := l.getUserIDFromContext()
	if userID == 0 {
		return nil, fmt.Errorf("未授权")
	}

	if err := r.ParseMultipartForm(propConfigMaxSize); err != nil {
		return nil, fmt.Errorf("解析表单失败: %w", err)
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("获取上传文件失败: %w", err)
	}
	defer file.Close()

	if fileHeader.Size > propConfigMaxSize {
		return nil, fmt.Errorf("文件大小超过限制")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	switch ext {
	case ".json", ".yaml", ".yml":
		// ok
	default:
		return nil, fmt.Errorf("不支持的文件类型")
	}

	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	dateDir := time.Now().Format("2006/01/02")
	storageDir := filepath.Join(l.svcCtx.Config.Upload.StoragePath, "prop-config", dateDir)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	dstPath := filepath.Join(storageDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer dst.Close()

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		buf := make([]byte, 512)
		n, _ := file.Read(buf)
		mimeType = http.DetectContentType(buf[:n])
		_, _ = file.Seek(0, 0)
	}

	if _, err := io.Copy(dst, file); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s/%s", l.svcCtx.Config.Upload.BaseURL, "prop-config", dateDir, filename)
	return &types.UploadPropConfigResp{
		Url:      url,
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		MimeType: mimeType,
	}, nil
}

func (l *UploadLogic) getUserIDFromContext() int64 {
	for _, key := range []string{"user_id", "userId"} {
		val := l.ctx.Value(key)
		if val == nil {
			continue
		}
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		default:
			continue
		}
	}
	return 0
}
