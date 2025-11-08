package user_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ur *UserRepository) HasRolePermission(ctx context.Context, roleID, permissionID uint) (bool, error) {
	var count int64
	err := ur.db.GetTx(ctx).WithContext(ctx).
		Model(&model.RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Count(&count).Error
	if err != nil {
		logger.Error(ctx, "Error checking role permission", "roleID", roleID, "permissionID", permissionID, err.Error())
		return false, err
	}
	return count > 0, err
}
