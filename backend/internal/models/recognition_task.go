package models

import (
	"database/sql"
	"time"
)

type RecognitionTask struct {
	ID           int64          `gorm:"primaryKey;column:id" json:"id"`
	TaskID       string         `gorm:"column:task_id;size:100;unique;not null" json:"task_id"`
	ImageID      int64          `gorm:"column:image_id;not null" json:"image_id"`
	UserID       int64          `gorm:"column:user_id;not null" json:"user_id"`
	ModelName    string         `gorm:"column:model_name;size:100" json:"model_name"`
	Status       int16          `gorm:"column:status;default:0" json:"status"` // 0-待处理 1-处理中 2-已完成 3-失败
	Progress     int            `gorm:"column:progress;default:0" json:"progress"`
	ResultCount  int            `gorm:"column:result_count;default:0" json:"result_count"`
	ErrorMessage string         `gorm:"column:error_message;type:text" json:"error_message,omitempty"`
	StartedAt    sql.NullTime   `gorm:"column:started_at" json:"started_at,omitempty"`
	CompletedAt  sql.NullTime   `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (RecognitionTask) TableName() string {
	return "recognition_tasks"
}

// 任务状态常量
const (
	TaskStatusPending    int16 = 0 // 待处理
	TaskStatusProcessing int16 = 1 // 处理中
	TaskStatusCompleted  int16 = 2 // 已完成
	TaskStatusFailed     int16 = 3 // 失败
)
