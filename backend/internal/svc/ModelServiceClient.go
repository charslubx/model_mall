package svc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// ModelServiceClient 模型服务客户端
type ModelServiceClient struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
}

// NewModelServiceClient 创建模型服务客户端
func NewModelServiceClient(baseURL, apiKey string) *ModelServiceClient {
	return &ModelServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

// RecognitionRequest 识别请求
type RecognitionRequest struct {
	TaskID    string `json:"task_id"`
	ImageURL  string `json:"image_url,omitempty"`
	ModelName string `json:"model_name,omitempty"`
	Callback  string `json:"callback,omitempty"` // 回调URL
}

// RecognitionResponse 识别响应
type RecognitionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TaskID string `json:"task_id"`
		Status string `json:"status"`
	} `json:"data"`
}

// UploadImageRequest 上传图片请求
type UploadImageRequest struct {
	TaskID    string
	FilePath  string
	ModelName string
	Callback  string
}

// UploadImageAndRecognize 上传图片并请求识别
func (c *ModelServiceClient) UploadImageAndRecognize(req *UploadImageRequest) (*RecognitionResponse, error) {
	// 打开文件
	file, err := os.Open(req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 创建multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件
	part, err := writer.CreateFormFile("image", req.FilePath)
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("复制文件内容失败: %w", err)
	}

	// 添加其他字段
	_ = writer.WriteField("task_id", req.TaskID)
	if req.ModelName != "" {
		_ = writer.WriteField("model_name", req.ModelName)
	}
	if req.Callback != "" {
		_ = writer.WriteField("callback", req.Callback)
	}

	writer.Close()

	// 创建请求
	url := fmt.Sprintf("%s/api/v1/recognize/upload", c.baseURL)
	httpReq, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var result RecognitionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// RecognizeByURL 通过URL请求识别
func (c *ModelServiceClient) RecognizeByURL(req *RecognitionRequest) (*RecognitionResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/recognize/url", c.baseURL)
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var result RecognitionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// TaskStatus 任务状态
type TaskStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		TaskID   string  `json:"task_id"`
		Status   string  `json:"status"`
		Progress int     `json:"progress"`
		Result   []Label `json:"result,omitempty"`
	} `json:"data"`
}

// Label 识别标签
type Label struct {
	Name       string                 `json:"name"`
	Code       string                 `json:"code,omitempty"`
	Confidence float64                `json:"confidence"`
	BBox       *BoundingBox           `json:"bbox,omitempty"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// BoundingBox 边界框
type BoundingBox struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// GetTaskStatus 获取任务状态
func (c *ModelServiceClient) GetTaskStatus(taskID string) (*TaskStatus, error) {
	url := fmt.Sprintf("%s/api/v1/task/%s/status", c.baseURL, taskID)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	var result TaskStatus
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}
