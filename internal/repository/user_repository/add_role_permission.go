package user_repository

import (
	"context"
	"wtm-backend/internal/infrastructure/database/model"
)

func (ur *UserRepository) AddRolePermission(ctx context.Context, roleID, permissionID uint) error {
	return ur.db.GetTx(ctx).WithContext(ctx).
		FirstOrCreate(&model.RolePermission{
			RoleID:       roleID,
			PermissionID: permissionID,
		}).Error
}
