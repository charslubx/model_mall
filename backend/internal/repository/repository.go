package repository

import (
	"gorm.io/gorm"
)

// Repositories 仓库集合
type Repositories struct {
	UserRepo                *UserRepository
	RoleRepo                *RoleRepository
	PermissionRepo          *PermissionRepository
	RolePermissionRepo      *RolePermissionRepository
	ImageRepo               *ImageRepository
	RecognitionTaskRepo     *RecognitionTaskRepository
	ClassificationLabelRepo *ClassificationLabelRepository
	ProductRepo             *ProductRepository
	CartRepo                *CartRepository
	OrderRepo               *OrderRepository
	AddressRepo             *AddressRepository
}

// NewRepositories 创建仓库集合
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:                NewUserRepository(db),
		RoleRepo:                NewRoleRepository(db),
		PermissionRepo:          NewPermissionRepository(db),
		RolePermissionRepo:      NewRolePermissionRepository(db),
		ImageRepo:               NewImageRepository(db),
		RecognitionTaskRepo:     NewRecognitionTaskRepository(db),
		ClassificationLabelRepo: NewClassificationLabelRepository(db),
		ProductRepo:             NewProductRepository(db),
		CartRepo:                NewCartRepository(db),
		OrderRepo:               NewOrderRepository(db),
		AddressRepo:             NewAddressRepository(db),
	}
}
