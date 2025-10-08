package models

import "time"

// Image 图片信息模型
type Image struct {
	ID         int64     `db:"id" json:"id"`
	Filename   string    `db:"filename" json:"filename"`         // 原始文件名
	FilePath   string    `db:"file_path" json:"file_path"`       // 文件存储路径
	FileSize   int64     `db:"file_size" json:"file_size"`       // 文件大小(字节)
	MimeType   string    `db:"mime_type" json:"mime_type"`       // 文件MIME类型
	Width      *int      `db:"width" json:"width,omitempty"`     // 图片宽度
	Height     *int      `db:"height" json:"height,omitempty"`   // 图片高度
	UploadedBy *int64    `db:"uploaded_by" json:"uploaded_by,omitempty"` // 上传用户ID
	Status     int       `db:"status" json:"status"`             // 状态：0-处理中 1-已分类 2-失败
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// ImageClassification 图片分类标签模型
type ImageClassification struct {
	ID           int64     `db:"id" json:"id"`
	ImageID      int64     `db:"image_id" json:"image_id"`           // 图片ID
	Label        string    `db:"label" json:"label"`                 // 分类标签
	Confidence   float64   `db:"confidence" json:"confidence"`       // 置信度(0-1)
	ModelName    string    `db:"model_name" json:"model_name"`       // 使用的模型名称
	ModelVersion *string   `db:"model_version" json:"model_version,omitempty"` // 模型版本
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// ImageStatus 图片处理状态常量
const (
	ImageStatusProcessing = 0 // 处理中
	ImageStatusClassified = 1 // 已分类
	ImageStatusFailed     = 2 // 失败
)

// ClassificationResult 分类结果
type ClassificationResult struct {
	Label      string  `json:"label"`      // 分类标签
	Confidence float64 `json:"confidence"` // 置信度
}