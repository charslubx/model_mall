package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"modelmall/backend/internal/models"
)

// ModelService 模型服务接口
type ModelService interface {
	// ClassifyImage 对图片进行分类
	ClassifyImage(ctx context.Context, imagePath string) ([]*models.ClassificationResult, error)
	// ClassifyImageFromBytes 对图片字节数据进行分类
	ClassifyImageFromBytes(ctx context.Context, imageData []byte, filename string) ([]*models.ClassificationResult, error)
}

// LocalModelService 本地模型服务实现
type LocalModelService struct {
	modelPath    string           // 模型文件路径
	modelName    string           // 模型名称
	modelVersion string           // 模型版本
	mu           sync.RWMutex     // 读写锁
	initialized  bool             // 是否已初始化
}

// RemoteModelService 远程模型服务实现（通过HTTP调用）
type RemoteModelService struct {
	endpoint     string // 远程模型服务端点
	modelName    string // 模型名称
	modelVersion string // 模型版本
	client       *http.Client
}

// NewLocalModelService 创建本地模型服务
func NewLocalModelService(modelPath, modelName, modelVersion string) *LocalModelService {
	return &LocalModelService{
		modelPath:    modelPath,
		modelName:    modelName,
		modelVersion: modelVersion,
		initialized:  false,
	}
}

// NewRemoteModelService 创建远程模型服务
func NewRemoteModelService(endpoint, modelName, modelVersion string) *RemoteModelService {
	return &RemoteModelService{
		endpoint:     endpoint,
		modelName:    modelName,
		modelVersion: modelVersion,
		client:       &http.Client{},
	}
}

// ClassifyImage 对图片进行分类（本地模型）
func (s *LocalModelService) ClassifyImage(ctx context.Context, imagePath string) ([]*models.ClassificationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 读取图片文件
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("读取图片文件失败: %w", err)
	}

	return s.ClassifyImageFromBytes(ctx, imageData, imagePath)
}

// ClassifyImageFromBytes 对图片字节数据进行分类（本地模型）
func (s *LocalModelService) ClassifyImageFromBytes(ctx context.Context, imageData []byte, filename string) ([]*models.ClassificationResult, error) {
	// 注意: 这里需要集成实际的机器学习模型
	// 示例使用模拟数据，实际应用中需要替换为真实的模型推理代码
	// 可以使用 TensorFlow Go、ONNX Runtime、或通过CGO调用Python模型等

	// 模拟分类结果
	results := []*models.ClassificationResult{
		{
			Label:      "cat",
			Confidence: 0.8523,
		},
		{
			Label:      "dog",
			Confidence: 0.1234,
		},
		{
			Label:      "bird",
			Confidence: 0.0243,
		},
	}

	return results, nil
}

// ClassifyImage 对图片进行分类（远程模型）
func (s *RemoteModelService) ClassifyImage(ctx context.Context, imagePath string) ([]*models.ClassificationResult, error) {
	// 读取图片文件
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		return nil, fmt.Errorf("读取图片文件失败: %w", err)
	}

	return s.ClassifyImageFromBytes(ctx, imageData, imagePath)
}

// ClassifyImageFromBytes 对图片字节数据进行分类（远程模型）
func (s *RemoteModelService) ClassifyImageFromBytes(ctx context.Context, imageData []byte, filename string) ([]*models.ClassificationResult, error) {
	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败: %w", err)
	}

	_, err = io.Copy(part, bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("写入文件数据失败: %w", err)
	}

	// 关闭writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭writer失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", s.endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("模型服务返回错误: %s, body: %s", resp.Status, string(bodyBytes))
	}

	// 解析响应
	var response struct {
		Results []*models.ClassificationResult `json:"results"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return response.Results, nil
}

// GetModelName 获取模型名称
func (s *LocalModelService) GetModelName() string {
	return s.modelName
}

// GetModelVersion 获取模型版本
func (s *LocalModelService) GetModelVersion() string {
	return s.modelVersion
}

// GetModelName 获取模型名称（远程）
func (s *RemoteModelService) GetModelName() string {
	return s.modelName
}

// GetModelVersion 获取模型版本（远程）
func (s *RemoteModelService) GetModelVersion() string {
	return s.modelVersion
}