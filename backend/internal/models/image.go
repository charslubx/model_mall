package models

import "time"

type Image struct {
	ID           int64     `gorm:"primaryKey;column:id" json:"id"`
	UserID       int64     `gorm:"column:user_id;not null" json:"user_id"`
	Filename     string    `gorm:"column:filename;size:255;not null" json:"filename"`
	OriginalName string    `gorm:"column:original_name;size:255;not null" json:"original_name"`
	FilePath     string    `gorm:"column:file_path;size:500;not null" json:"file_path"`
	FileURL      string    `gorm:"column:file_url;size:500" json:"file_url"`
	FileSize     int64     `gorm:"column:file_size;not null" json:"file_size"`
	MimeType     string    `gorm:"column:mime_type;size:100;not null" json:"mime_type"`
	Width        int       `gorm:"column:width" json:"width"`
	Height       int       `gorm:"column:height" json:"height"`
	MD5          string    `gorm:"column:md5;size:32" json:"md5"`
	Status       int16     `gorm:"column:status;default:1" json:"status"` // 0-已删除 1-正常
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Image) TableName() string {
	return "images"
}
