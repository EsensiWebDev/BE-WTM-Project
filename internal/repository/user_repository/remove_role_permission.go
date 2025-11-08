package user_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
)

func (ur *UserRepository) RemoveRolePermission(ctx context.Context, roleID, permissionID uint) error {
	return ur.db.GetTx(ctx).WithContext(ctx).
		Unscoped().
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&model.RolePermission{}).Error
}
