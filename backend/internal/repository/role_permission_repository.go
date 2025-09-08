package repository

import (
	"context"

	"model_mall_backend/backend/internal/models"
	"model_mall_backend/backend/internal/svc"
	
	"gorm.io/gorm"
)

type RolePermissionRepository struct {
	svcCtx *svc.ServiceContext
}

func NewRolePermissionRepository(svcCtx *svc.ServiceContext) *RolePermissionRepository {
	return &RolePermissionRepository{
		svcCtx: svcCtx,
	}
}

// BatchAssign 批量分配角色权限
func (r *RolePermissionRepository) BatchAssign(ctx context.Context, roleID int64, permissionIDs []int64) error {
	// 开启事务
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除原有权限
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// 批量插入新权限
		if len(permissionIDs) > 0 {
			var rolePermissions []models.RolePermission
			for _, permissionID := range permissionIDs {
				rolePermissions = append(rolePermissions, models.RolePermission{
					RoleID:       roleID,
					PermissionID: permissionID,
				})
			}
			
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetByRoleID 根据角色ID获取权限关联
func (r *RolePermissionRepository) GetByRoleID(ctx context.Context, roleID int64) ([]*models.RolePermission, error) {
	var rolePermissions []*models.RolePermission
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Preload("Permission").
		Where("role_id = ?", roleID).
		Find(&rolePermissions).Error
	return rolePermissions, err
}

// GetPermissionIDsByRoleID 根据角色ID获取权限ID列表
func (r *RolePermissionRepository) GetPermissionIDsByRoleID(ctx context.Context, roleID int64) ([]int64, error) {
	var permissionIDs []int64
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Model(&models.RolePermission{}).
		Where("role_id = ?", roleID).
		Pluck("permission_id", &permissionIDs).Error
	return permissionIDs, err
}

// DeleteByRoleID 根据角色ID删除权限关联
func (r *RolePermissionRepository) DeleteByRoleID(ctx context.Context, roleID int64) error {
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Where("role_id = ?", roleID).
		Delete(&models.RolePermission{}).Error
}

// DeleteByPermissionID 根据权限ID删除权限关联
func (r *RolePermissionRepository) DeleteByPermissionID(ctx context.Context, permissionID int64) error {
	return r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Where("permission_id = ?", permissionID).
		Delete(&models.RolePermission{}).Error
}

// HasPermission 检查角色是否有指定权限
func (r *RolePermissionRepository) HasPermission(ctx context.Context, roleID int64, permissionCode string) (bool, error) {
	var count int64
	err := r.svcCtx.OrmHelper.GetDB().WithContext(ctx).
		Model(&models.RolePermission{}).
		Joins("JOIN permissions p ON p.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND p.code = ? AND p.status = ?", roleID, permissionCode, 1).
		Count(&count).Error
	return count > 0, err
}