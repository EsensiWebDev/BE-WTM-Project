package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetAllRolesWithPermissions(ctx context.Context) ([]entity.Role, error) {
	db := ur.db.GetTx(ctx)

	var roles []model.Role
	err := db.WithContext(ctx).Debug().
		Where("id != ?", constant.DefaultRoleSuperAdmin).
		Preload("Permissions").
		Find(&roles).Error
	if err != nil {
		if ur.db.ErrRecordNotFound(ctx, err) {
			logger.Warn(ctx, "No roles found in the database")
			return nil, nil
		}
		logger.Error(ctx, "Error retrieving all roles with permissions", err.Error())
		return nil, err
	}

	var entityRoles []entity.Role
	if err := utils.CopyPatch(&entityRoles, &roles); err != nil {
		logger.Error(ctx, "Error copying roles model to entity", err.Error())
		return nil, err
	}
	return entityRoles, nil
}
