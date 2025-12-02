package user_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
)

func (ur *UserRepository) AddRolePermission(ctx context.Context, roleID, permissionID uint) error {
	if err := ur.db.GetTx(ctx).WithContext(ctx).
		Debug().
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		FirstOrCreate(&model.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}).Error; err != nil {
		logger.Error(ctx, "Error when create role permission", err.Error())
		return err
	}
	return nil
}
