package repository

import (
	"context"
	"fmt"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
)

type PermissionRepository struct {
	svcCtx *svc.ServiceContext
}

func NewPermissionRepository(svcCtx *svc.ServiceContext) *PermissionRepository {
	return &PermissionRepository{
		svcCtx: svcCtx,
	}
}

// Create 创建权限
func (r *PermissionRepository) Create(ctx context.Context, permission *models.Permission) error {
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Create(permission).Error
}

// GetByID 根据ID获取权限
func (r *PermissionRepository) GetByID(ctx context.Context, id int64) (*models.Permission, error) {
	var permission models.Permission
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Preload("Children").
		Preload("Parent").
		First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetByCode 根据代码获取权限
func (r *PermissionRepository) GetByCode(ctx context.Context, code string) (*models.Permission, error) {
	var permission models.Permission
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Where("code = ?", code).
		First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// Update 更新权限
func (r *PermissionRepository) Update(ctx context.Context, permission *models.Permission) error {
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Save(permission).Error
}

// UpdateByID 根据ID更新权限
func (r *PermissionRepository) UpdateByID(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Model(&models.Permission{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// Delete 删除权限
func (r *PermissionRepository) Delete(ctx context.Context, id int64) error {
	// 检查是否为系统权限
	var permission models.Permission
	if err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).First(&permission, id).Error; err != nil {
		return err
	}
	if permission.IsSystem {
		return fmt.Errorf("不能删除系统权限")
	}

	// 检查是否有子权限
	var childCount int64
	if err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Model(&models.Permission{}).
		Where("parent_id = ?", id).
		Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return fmt.Errorf("该权限有 %d 个子权限，请先删除子权限", childCount)
	}

	// 删除角色权限关联
	if err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Where("permission_id = ?", id).
		Delete(&models.RolePermission{}).Error; err != nil {
		return err
	}

	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Delete(&models.Permission{}, id).Error
}

// List 获取权限列表
func (r *PermissionRepository) List(ctx context.Context, req *models.PermissionListReq) ([]*models.Permission, int64, error) {
	db := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Model(&models.Permission{})

	// 构建查询条件
	if req.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Type != "" {
		db = db.Where("type = ?", req.Type)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.ParentID != nil {
		db = db.Where("parent_id = ?", *req.ParentID)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var permissions []*models.Permission
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := db.Preload("Children").
		Preload("Parent").
		Offset(offset).
		Limit(pageSize).
		Order("sort ASC, created_at DESC").
		Find(&permissions).Error

	return permissions, total, err
}

// GetTree 获取权限树形结构
func (r *PermissionRepository) GetTree(ctx context.Context) ([]*models.PermissionTree, error) {
	var permissions []*models.Permission
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Where("status = ?", 1).
		Order("sort ASC, created_at DESC").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}

	return r.buildPermissionTree(permissions, 0), nil
}

// buildPermissionTree 构建权限树
func (r *PermissionRepository) buildPermissionTree(permissions []*models.Permission, parentID int64) []*models.PermissionTree {
	var tree []*models.PermissionTree
	
	for _, permission := range permissions {
		if permission.ParentID == parentID {
			node := &models.PermissionTree{
				ID:       permission.ID,
				Name:     permission.Name,
				Code:     permission.Code,
				Type:     permission.Type,
				ParentID: permission.ParentID,
				Sort:     permission.Sort,
			}
			
			// 递归获取子权限
			children := r.buildPermissionTree(permissions, permission.ID)
			if len(children) > 0 {
				node.Children = children
			}
			
			tree = append(tree, node)
		}
	}
	
	return tree
}

// GetByRoleID 根据角色ID获取权限列表
func (r *PermissionRepository) GetByRoleID(ctx context.Context, roleID int64) ([]*models.Permission, error) {
	var permissions []*models.Permission
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Joins("JOIN role_permissions rp ON rp.permission_id = permissions.id").
		Where("rp.role_id = ? AND permissions.status = ?", roleID, 1).
		Order("permissions.sort ASC").
		Find(&permissions).Error
	return permissions, err
}

// ExistsByCode 检查权限代码是否存在
func (r *PermissionRepository) ExistsByCode(ctx context.Context, code string, excludeID ...int64) (bool, error) {
	db := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Model(&models.Permission{})
	db = db.Where("code = ?", code)
	
	if len(excludeID) > 0 && excludeID[0] > 0 {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}