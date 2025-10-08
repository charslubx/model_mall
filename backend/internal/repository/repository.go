package repository

import (
	"gorm.io/gorm"
)

// Repositories 仓库集合
type Repositories struct {
	User                  *UserRepository
	Role                  *RoleRepository
	Permission            *PermissionRepository
	RolePermission        *RolePermissionRepository
	Image                 *ImageRepository
	RecognitionTask       *RecognitionTaskRepository
	ClassificationLabel   *ClassificationLabelRepository
}

// NewRepositories 创建仓库集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:                NewUserRepository(db),
		Role:                NewRoleRepository(db),
		Permission:          NewPermissionRepository(db),
		RolePermission:      NewRolePermissionRepository(db),
		Image:               NewImageRepository(db),
		RecognitionTask:     NewRecognitionTaskRepository(db),
		ClassificationLabel: NewClassificationLabelRepository(db),
	}
}