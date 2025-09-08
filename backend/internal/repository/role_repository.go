package repository

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
)

type RoleRepository struct {
	svcCtx *svc.ServiceContext
}

func NewRoleRepository(svcCtx *svc.ServiceContext) *RoleRepository {
	return &RoleRepository{
		svcCtx: svcCtx,
	}
}

// Create 创建角色
func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Create(role).Error
}

// GetByID 根据ID获取角色
func (r *RoleRepository) GetByID(ctx context.Context, id int64) (*models.Role, error) {
	var role models.Role
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Preload("Permissions").
		First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据代码获取角色
func (r *RoleRepository) GetByCode(ctx context.Context, code string) (*models.Role, error) {
	var role models.Role
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Preload("Permissions").
		Where("code = ?", code).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// Update 更新角色
func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Save(role).Error
}

// UpdateByID 根据ID更新角色
func (r *RoleRepository) UpdateByID(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Model(&models.Role{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(ctx context.Context, id int64) error {
	// 检查是否为系统角色
	var role models.Role
	if err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).First(&role, id).Error; err != nil {
		return err
	}
	if role.IsSystem {
		return fmt.Errorf("不能删除系统角色")
	}

	// 检查是否有用户使用该角色
	var userCount int64
	if err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Model(&models.User{}).
		Where("role_id = ?", id).
		Count(&userCount).Error; err != nil {
		return err
	}
	if userCount > 0 {
		return fmt.Errorf("该角色已被 %d 个用户使用，无法删除", userCount)
	}

	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Delete(&models.Role{}, id).Error
}

// List 获取角色列表
func (r *RoleRepository) List(ctx context.Context, req *models.RoleListReq) ([]*models.Role, int64, error) {
	db := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Model(&models.Role{})

	// 构建查询条件
	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var roles []*models.Role
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := db.Preload("Permissions").
		Offset(offset).
		Limit(pageSize).
		Order("sort ASC, created_at DESC").
		Find(&roles).Error

	return roles, total, err
}

// GetAll 获取所有角色（用于选项列表）
func (r *RoleRepository) GetAll(ctx context.Context) ([]*models.RoleOption, error) {
	var roles []*models.RoleOption
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Model(&models.Role{}).
		Select("id, name, code").
		Where("status = ?", 1).
		Order("sort ASC, created_at DESC").
		Find(&roles).Error
	return roles, err
}

// ExistsByName 检查角色名是否存在
func (r *RoleRepository) ExistsByName(ctx context.Context, name string, excludeID ...int64) (bool, error) {
	db := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Model(&models.Role{})
	db = db.Where("name = ?", name)
	
	if len(excludeID) > 0 && excludeID[0] > 0 {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}

// ExistsByCode 检查角色代码是否存在
func (r *RoleRepository) ExistsByCode(ctx context.Context, code string, excludeID ...int64) (bool, error) {
	db := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Model(&models.Role{})
	db = db.Where("code = ?", code)
	
	if len(excludeID) > 0 && excludeID[0] > 0 {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}