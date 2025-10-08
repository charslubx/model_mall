package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	"model_mall_backend/backend/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type ImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ImageLogic {
	return &ImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UploadImage 上传图片
func (l *ImageLogic) UploadImage(r *http.Request, req *types.UploadImageReq) (*types.UploadImageResp, error) {
	// 获取用户ID（从context中获取，这里假设从JWT中间件已经设置）
	userID := l.getUserIDFromContext()
	if userID == 0 {
		return nil, fmt.Errorf("未授权")
	}

	// 解析multipart form
	err := r.ParseMultipartForm(l.svcCtx.Config.Upload.MaxSize)
	if err != nil {
		return nil, fmt.Errorf("解析表单失败: %w", err)
	}

	// 获取上传的文件
	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		return nil, fmt.Errorf("获取上传文件失败: %w", err)
	}
	defer file.Close()

	// 验证文件类型
	if !l.isAllowedFileType(fileHeader.Filename) {
		return nil, fmt.Errorf("不支持的文件类型")
	}

	// 验证文件大小
	if fileHeader.Size > l.svcCtx.Config.Upload.MaxSize {
		return nil, fmt.Errorf("文件大小超过限制")
	}

	// 计算文件MD5
	md5Hash, err := l.calculateMD5(file)
	if err != nil {
		return nil, fmt.Errorf("计算文件MD5失败: %w", err)
	}
	file.Seek(0, 0) // 重置文件指针

	// 检查文件是否已存在
	existingImage, _ := l.svcCtx.Repos.Image.GetByMD5(md5Hash)
	if existingImage != nil {
		// 文件已存在，返回已有记录
		return &types.UploadImageResp{
			ImageID:  existingImage.ID,
			TaskID:   "",
			FileURL:  existingImage.FileURL,
			Filename: existingImage.Filename,
		}, nil
	}

	// 生成文件名
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	
	// 按日期组织存储路径
	dateDir := time.Now().Format("2006/01/02")
	storageDir := filepath.Join(l.svcCtx.Config.Upload.StoragePath, dateDir)
	
	// 确保目录存在
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %w", err)
	}

	// 保存文件
	filePath := filepath.Join(storageDir, filename)
	destFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, file); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 构建访问URL
	fileURL := fmt.Sprintf("%s/%s/%s", l.svcCtx.Config.Upload.BaseURL, dateDir, filename)

	// 保存图片记录
	image := &models.Image{
		UserID:       userID,
		Filename:     filename,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileURL:      fileURL,
		FileSize:     fileHeader.Size,
		MimeType:     fileHeader.Header.Get("Content-Type"),
		MD5:          md5Hash,
		Status:       1,
	}

	if err := l.svcCtx.Repos.Image.Create(image); err != nil {
		return nil, fmt.Errorf("保存图片记录失败: %w", err)
	}

	// 创建识别任务
	taskID := uuid.New().String()
	task := &models.RecognitionTask{
		TaskID:    taskID,
		ImageID:   image.ID,
		UserID:    userID,
		ModelName: req.ModelName,
		Status:    models.TaskStatusPending,
		Progress:  0,
	}

	if err := l.svcCtx.Repos.RecognitionTask.Create(task); err != nil {
		return nil, fmt.Errorf("创建识别任务失败: %w", err)
	}

	// 异步调用模型服务
	go l.callModelService(taskID, filePath, fileURL, req.ModelName)

	return &types.UploadImageResp{
		ImageID:  image.ID,
		TaskID:   taskID,
		FileURL:  fileURL,
		Filename: filename,
	}, nil
}

// GetImageList 获取图片列表
func (l *ImageLogic) GetImageList(req *types.GetImageListReq) (*types.GetImageListResp, error) {
	userID := l.getUserIDFromContext()
	if userID == 0 {
		return nil, fmt.Errorf("未授权")
	}

	images, total, err := l.svcCtx.Repos.Image.GetByUserID(userID, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取图片列表失败: %w", err)
	}

	list := make([]types.ImageInfo, 0, len(images))
	for _, img := range images {
		list = append(list, types.ImageInfo{
			ID:           img.ID,
			Filename:     img.Filename,
			OriginalName: img.OriginalName,
			FileURL:      img.FileURL,
			FileSize:     img.FileSize,
			Width:        img.Width,
			Height:       img.Height,
			Status:       img.Status,
			CreatedAt:    img.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &types.GetImageListResp{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}, nil
}

// 辅助方法

func (l *ImageLogic) getUserIDFromContext() int64 {
	// 从context中获取用户ID
	// 这里需要根据实际的JWT中间件实现来获取
	// 示例：
	userIDVal := l.ctx.Value("user_id")
	if userIDVal == nil {
		return 0
	}
	
	if userID, ok := userIDVal.(int64); ok {
		return userID
	}
	
	return 0
}

func (l *ImageLogic) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedTypes := strings.Split(l.svcCtx.Config.Upload.AllowedTypes, ",")
	
	for _, t := range allowedTypes {
		if ext == strings.TrimSpace(t) {
			return true
		}
	}
	
	return false
}

func (l *ImageLogic) calculateMD5(file io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (l *ImageLogic) callModelService(taskID, filePath, fileURL, modelName string) {
	// 更新任务状态为处理中
	_ = l.svcCtx.Repos.RecognitionTask.UpdateStatus(taskID, models.TaskStatusProcessing)

	// 构建回调URL
	callbackURL := fmt.Sprintf("%s/api/model/callback", l.svcCtx.Config.Host)

	// 调用模型服务
	_, err := l.svcCtx.ModelServiceClient.UploadImageAndRecognize(&svc.UploadImageRequest{
		TaskID:    taskID,
		FilePath:  filePath,
		ModelName: modelName,
		Callback:  callbackURL,
	})

	if err != nil {
		logx.Errorf("调用模型服务失败: %v", err)
		_ = l.svcCtx.Repos.RecognitionTask.UpdateError(taskID, err.Error())
	}
}
