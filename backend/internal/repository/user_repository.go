package repository

import (
	"context"

	"model_mall_backend/backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create 创建用户
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Role.Permissions").
		First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("Role.Permissions").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateByID 根据ID更新用户
func (r *UserRepository) UpdateByID(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// List 获取用户列表
func (r *UserRepository) List(ctx context.Context, req *models.UserListReq) ([]*models.User, int64, error) {
	db := r.db.WithContext(ctx).Model(&models.User{})

	// 构建查询条件
	if req.Keyword != "" {
		db = db.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.RoleID != nil {
		db = db.Where("role_id = ?", *req.RoleID)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var users []*models.User
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := db.Preload("Role").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string, excludeID ...int64) (bool, error) {
	db := r.db.WithContext(ctx).Model(&models.User{})
	db = db.Where("username = ?", username)
	
	if len(excludeID) > 0 && excludeID[0] > 0 {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string, excludeID ...int64) (bool, error) {
	db := r.db.WithContext(ctx).Model(&models.User{})
	db = db.Where("email = ?", email)
	
	if len(excludeID) > 0 && excludeID[0] > 0 {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}

// ExistsByPhone 检查手机号是否存在
func (r *UserRepository) ExistsByPhone(ctx context.Context, phone string, excludeID ...int64) (bool, error) {
	if phone == "" {
		return false, nil
	}

	db := r.db.WithContext(ctx).Model(&models.User{})
	db = db.Where("phone = ?", phone)
	
	if len(excludeID) > 0 && excludeID[0] > 0 {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}

// UpdateLastLogin 更新最后登录信息
func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int64, ip string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": gorm.Expr("NOW()"),
			"last_login_ip": ip,
		}).Error
}