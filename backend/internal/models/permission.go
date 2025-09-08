package models

import (
	"time"
)

// Permission 权限表
type Permission struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:权限ID"`
	Name        string    `json:"name" gorm:"type:varchar(50);not null;comment:权限名称"`
	Code        string    `json:"code" gorm:"type:varchar(100);uniqueIndex;not null;comment:权限代码"`
	Type        string    `json:"type" gorm:"type:varchar(20);not null;comment:权限类型 menu-菜单 button-按钮 api-接口"`
	ParentID    int64     `json:"parent_id" gorm:"default:0;comment:父权限ID"`
	Path        string    `json:"path" gorm:"type:varchar(255);comment:路径/接口地址"`
	Method      string    `json:"method" gorm:"type:varchar(10);comment:请求方法"`
	Icon        string    `json:"icon" gorm:"type:varchar(100);comment:图标"`
	Component   string    `json:"component" gorm:"type:varchar(255);comment:组件路径"`
	Sort        int       `json:"sort" gorm:"default:0;comment:排序"`
	Status      int8      `json:"status" gorm:"type:tinyint;default:1;comment:状态 0-禁用 1-正常"`
	IsSystem    bool      `json:"is_system" gorm:"default:false;comment:是否系统权限"`
	Description string    `json:"description" gorm:"type:varchar(255);comment:权限描述"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	
	// 关联关系
	Children []Permission `json:"children" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Parent   *Permission  `json:"parent" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// PermissionCreateReq 创建权限请求
type PermissionCreateReq struct {
	Name        string `json:"name" validate:"required,min=2,max=50"`
	Code        string `json:"code" validate:"required,min=2,max=100"`
	Type        string `json:"type" validate:"required,oneof=menu button api"`
	ParentID    int64  `json:"parent_id" validate:"omitempty,min=0"`
	Path        string `json:"path" validate:"omitempty,max=255"`
	Method      string `json:"method" validate:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	Icon        string `json:"icon" validate:"omitempty,max=100"`
	Component   string `json:"component" validate:"omitempty,max=255"`
	Sort        int    `json:"sort" validate:"omitempty,min=0"`
	Status      int8   `json:"status" validate:"omitempty,oneof=0 1"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

// PermissionUpdateReq 更新权限请求
type PermissionUpdateReq struct {
	ID          int64  `json:"id" validate:"required"`
	Name        string `json:"name" validate:"omitempty,min=2,max=50"`
	Code        string `json:"code" validate:"omitempty,min=2,max=100"`
	Type        string `json:"type" validate:"omitempty,oneof=menu button api"`
	ParentID    int64  `json:"parent_id" validate:"omitempty,min=0"`
	Path        string `json:"path" validate:"omitempty,max=255"`
	Method      string `json:"method" validate:"omitempty,oneof=GET POST PUT DELETE PATCH"`
	Icon        string `json:"icon" validate:"omitempty,max=100"`
	Component   string `json:"component" validate:"omitempty,max=255"`
	Sort        int    `json:"sort" validate:"omitempty,min=0"`
	Status      int8   `json:"status" validate:"omitempty,oneof=0 1"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

// PermissionResp 权限响应
type PermissionResp struct {
	ID          int64              `json:"id"`
	Name        string             `json:"name"`
	Code        string             `json:"code"`
	Type        string             `json:"type"`
	ParentID    int64              `json:"parent_id"`
	Path        string             `json:"path"`
	Method      string             `json:"method"`
	Icon        string             `json:"icon"`
	Component   string             `json:"component"`
	Sort        int                `json:"sort"`
	Status      int8               `json:"status"`
	IsSystem    bool               `json:"is_system"`
	Description string             `json:"description"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Children    []PermissionResp   `json:"children,omitempty"`
}

// PermissionListReq 权限列表请求
type PermissionListReq struct {
	Page     int    `json:"page" validate:"omitempty,min=1"`
	PageSize int    `json:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword  string `json:"keyword" validate:"omitempty,max=50"`
	Type     string `json:"type" validate:"omitempty,oneof=menu button api"`
	Status   *int8  `json:"status" validate:"omitempty,oneof=0 1"`
	ParentID *int64 `json:"parent_id" validate:"omitempty,min=0"`
}

// PermissionTree 权限树形结构
type PermissionTree struct {
	ID       int64            `json:"id"`
	Name     string           `json:"name"`
	Code     string           `json:"code"`
	Type     string           `json:"type"`
	ParentID int64            `json:"parent_id"`
	Sort     int              `json:"sort"`
	Children []PermissionTree `json:"children,omitempty"`
}