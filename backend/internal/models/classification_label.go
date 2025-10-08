package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONB 自定义类型用于PostgreSQL的JSONB字段
type JSONB map[string]interface{}

// Value 实现 driver.Valuer 接口
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现 sql.Scanner 接口
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type ClassificationLabel struct {
	ID         int64     `gorm:"primaryKey;column:id" json:"id"`
	TaskID     int64     `gorm:"column:task_id;not null" json:"task_id"`
	ImageID    int64     `gorm:"column:image_id;not null" json:"image_id"`
	LabelName  string    `gorm:"column:label_name;size:100;not null" json:"label_name"`
	LabelCode  string    `gorm:"column:label_code;size:100" json:"label_code"`
	Confidence float64   `gorm:"column:confidence;type:decimal(5,4)" json:"confidence"` // 0.0000-1.0000
	BboxX      int       `gorm:"column:bbox_x" json:"bbox_x,omitempty"`
	BboxY      int       `gorm:"column:bbox_y" json:"bbox_y,omitempty"`
	BboxWidth  int       `gorm:"column:bbox_width" json:"bbox_width,omitempty"`
	BboxHeight int       `gorm:"column:bbox_height" json:"bbox_height,omitempty"`
	ExtraData  JSONB     `gorm:"column:extra_data;type:jsonb" json:"extra_data,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (ClassificationLabel) TableName() string {
	return "classification_labels"
}
