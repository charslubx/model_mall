package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"model_mall_backend/backend/internal/models"
)

// ModelService 模型服务
type ModelService struct {
	modelEndpoint string
	uploadPath    string
	httpClient    *http.Client
}

// NewModelService 创建模型服务实例
func NewModelService(modelEndpoint, uploadPath string) *ModelService {
	return &ModelService{
		modelEndpoint: modelEndpoint,
		uploadPath:    uploadPath,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProcessImage 处理图片分类
func (s *ModelService) ProcessImage(ctx context.Context, req *models.ImageClassificationReq) (*models.ModelResponse, error) {
	startTime := time.Now()
	
	// 创建multipart表单数据
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// 添加图片文件
	part, err := writer.CreateFormFile("image", req.ImageName)
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败: %v", err)
	}
	
	_, err = part.Write(req.ImageData)
	if err != nil {
		return nil, fmt.Errorf("写入图片数据失败: %v", err)
	}
	
	// 添加模型名称
	err = writer.WriteField("model_name", req.ModelName)
	if err != nil {
		return nil, fmt.Errorf("写入模型名称失败: %v", err)
	}
	
	// 添加最小置信度
	if req.MinConfidence > 0 {
		err = writer.WriteField("min_confidence", fmt.Sprintf("%.4f", req.MinConfidence))
		if err != nil {
			return nil, fmt.Errorf("写入最小置信度失败: %v", err)
		}
	}
	
	writer.Close()
	
	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.modelEndpoint+"/predict", &buf)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	
	// 发送请求
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送模型请求失败: %v", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("模型服务返回错误: %d, %s", resp.StatusCode, string(respBody))
	}
	
	// 解析响应
	var modelResp models.ModelResponse
	err = json.Unmarshal(respBody, &modelResp)
	if err != nil {
		return nil, fmt.Errorf("解析模型响应失败: %v", err)
	}
	
	// 设置处理时间
	modelResp.ProcessTime = time.Since(startTime).Milliseconds()
	
	return &modelResp, nil
}

// SaveImage 保存图片到本地
func (s *ModelService) SaveImage(imageData []byte, imageName string) (string, error) {
	// 确保上传目录存在
	err := os.MkdirAll(s.uploadPath, 0755)
	if err != nil {
		return "", fmt.Errorf("创建上传目录失败: %v", err)
	}
	
	// 生成唯一文件名
	ext := filepath.Ext(imageName)
	hash := fmt.Sprintf("%x", md5.Sum(imageData))
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s_%s%s", timestamp, hash[:8], ext)
	
	// 创建日期子目录
	dateDir := time.Now().Format("2006/01/02")
	fullDir := filepath.Join(s.uploadPath, dateDir)
	err = os.MkdirAll(fullDir, 0755)
	if err != nil {
		return "", fmt.Errorf("创建日期目录失败: %v", err)
	}
	
	// 完整文件路径
	filePath := filepath.Join(fullDir, fileName)
	
	// 写入文件
	err = os.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return "", fmt.Errorf("保存图片失败: %v", err)
	}
	
	// 返回相对路径
	relativePath := filepath.Join(dateDir, fileName)
	return relativePath, nil
}

// GetImageFormat 获取图片格式
func (s *ModelService) GetImageFormat(imageData []byte) string {
	if len(imageData) < 4 {
		return "unknown"
	}
	
	// 检测常见图片格式
	switch {
	case bytes.HasPrefix(imageData, []byte{0xFF, 0xD8, 0xFF}):
		return "jpeg"
	case bytes.HasPrefix(imageData, []byte{0x89, 0x50, 0x4E, 0x47}):
		return "png"
	case bytes.HasPrefix(imageData, []byte{0x47, 0x49, 0x46}):
		return "gif"
	case bytes.HasPrefix(imageData, []byte{0x52, 0x49, 0x46, 0x46}) && 
		 bytes.Contains(imageData[:12], []byte("WEBP")):
		return "webp"
	case bytes.HasPrefix(imageData, []byte{0x42, 0x4D}):
		return "bmp"
	default:
		return "unknown"
	}
}

// ValidateImage 验证图片数据
func (s *ModelService) ValidateImage(imageData []byte, maxSize int64) error {
	if len(imageData) == 0 {
		return fmt.Errorf("图片数据为空")
	}
	
	if int64(len(imageData)) > maxSize {
		return fmt.Errorf("图片大小超过限制: %d bytes", maxSize)
	}
	
	format := s.GetImageFormat(imageData)
	if format == "unknown" {
		return fmt.Errorf("不支持的图片格式")
	}
	
	return nil
}

// GetSupportedFormats 获取支持的图片格式
func (s *ModelService) GetSupportedFormats() []string {
	return []string{"jpeg", "jpg", "png", "gif", "webp", "bmp"}
}

// IsFormatSupported 检查格式是否支持
func (s *ModelService) IsFormatSupported(format string) bool {
	format = strings.ToLower(format)
	if format == "jpg" {
		format = "jpeg"
	}
	
	supportedFormats := s.GetSupportedFormats()
	for _, supported := range supportedFormats {
		if format == supported {
			return true
		}
	}
	return false
}

// HealthCheck 模型服务健康检查
func (s *ModelService) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", s.modelEndpoint+"/health", nil)
	if err != nil {
		return fmt.Errorf("创建健康检查请求失败: %v", err)
	}
	
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("模型服务健康检查失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("模型服务不健康: %d", resp.StatusCode)
	}
	
	return nil
}