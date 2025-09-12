package models

import (
	"time"
)

// ImageClassification 图片分类记录表
type ImageClassification struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:分类记录ID"`
	ImagePath   string    `json:"image_path" gorm:"type:varchar(500);not null;comment:图片路径"`
	ImageName   string    `json:"image_name" gorm:"type:varchar(255);not null;comment:图片名称"`
	ImageSize   int64     `json:"image_size" gorm:"comment:图片大小(字节)"`
	ImageFormat string    `json:"image_format" gorm:"type:varchar(20);comment:图片格式"`
	ModelName   string    `json:"model_name" gorm:"type:varchar(100);not null;comment:使用的模型名称"`
	ModelVersion string   `json:"model_version" gorm:"type:varchar(50);comment:模型版本"`
	ProcessTime int64     `json:"process_time" gorm:"comment:处理耗时(毫秒)"`
	Confidence  float64   `json:"confidence" gorm:"type:decimal(5,4);comment:总体置信度"`
	Status      int8      `json:"status" gorm:"type:smallint;default:1;comment:状态 0-失败 1-成功 2-处理中"`
	UserID      int64     `json:"user_id" gorm:"comment:用户ID"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	
	// 关联关系
	User   User                     `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Labels []ImageClassificationLabel `json:"labels" gorm:"foreignKey:ClassificationID"`
}

// TableName 指定表名
func (ImageClassification) TableName() string {
	return "image_classifications"
}

// ImageClassificationLabel 图片分类标签表
type ImageClassificationLabel struct {
	ID               int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:标签ID"`
	ClassificationID int64     `json:"classification_id" gorm:"not null;comment:分类记录ID"`
	LabelName        string    `json:"label_name" gorm:"type:varchar(100);not null;comment:标签名称"`
	LabelCode        string    `json:"label_code" gorm:"type:varchar(100);comment:标签代码"`
	Confidence       float64   `json:"confidence" gorm:"type:decimal(5,4);not null;comment:置信度"`
	BoundingBox      string    `json:"bounding_box" gorm:"type:text;comment:边界框信息(JSON格式)"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	
	// 关联关系
	Classification ImageClassification `json:"classification" gorm:"foreignKey:ClassificationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName 指定表名
func (ImageClassificationLabel) TableName() string {
	return "image_classification_labels"
}

// ImageClassificationReq 图片分类请求
type ImageClassificationReq struct {
	ImageData   []byte `json:"image_data" validate:"required"`
	ImageName   string `json:"image_name" validate:"required"`
	ModelName   string `json:"model_name" validate:"required"`
	UserID      int64  `json:"user_id"`
	SaveImage   bool   `json:"save_image"`
	MinConfidence float64 `json:"min_confidence" validate:"omitempty,min=0,max=1"`
}

// ImageClassificationResp 图片分类响应
type ImageClassificationResp struct {
	ID          int64                     `json:"id"`
	ImagePath   string                    `json:"image_path"`
	ImageName   string                    `json:"image_name"`
	ModelName   string                    `json:"model_name"`
	ProcessTime int64                     `json:"process_time"`
	Confidence  float64                   `json:"confidence"`
	Status      int8                      `json:"status"`
	Labels      []ImageClassificationLabel `json:"labels"`
	CreatedAt   time.Time                 `json:"created_at"`
}

// BoundingBox 边界框结构
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ModelPrediction 模型预测结果
type ModelPrediction struct {
	Label       string       `json:"label"`
	Code        string       `json:"code"`
	Confidence  float64      `json:"confidence"`
	BoundingBox *BoundingBox `json:"bounding_box,omitempty"`
}

// ModelResponse 模型服务响应
type ModelResponse struct {
	Success     bool              `json:"success"`
	Message     string            `json:"message"`
	ProcessTime int64             `json:"process_time"`
	Predictions []ModelPrediction `json:"predictions"`
}