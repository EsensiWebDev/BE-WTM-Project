package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetAllPermissions(ctx context.Context) ([]entity.Permission, error) {
	var perms []model.Permission
	err := ur.db.GetTx(ctx).WithContext(ctx).Find(&perms).Error
	if err != nil {
		return nil, err
	}

	var entityPerms []entity.Permission
	if err := utils.CopyPatch(&entityPerms, &perms); err != nil {
		logger.Error(ctx, "Error copying permissions model to entity", err.Error())
		return nil, err
	}
	return entityPerms, nil
}
