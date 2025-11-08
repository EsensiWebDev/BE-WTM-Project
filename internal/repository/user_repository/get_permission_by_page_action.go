package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetPermissionByPageAction(ctx context.Context, page, action string) (*entity.Permission, error) {
	var perm model.Permission
	err := ur.db.GetTx(ctx).WithContext(ctx).
		Where("page = ? AND action = ?", page, action).
		First(&perm).Error
	if err != nil {
		return nil, err
	}

	var entityPerm entity.Permission
	if err := utils.CopyPatch(&entityPerm, &perm); err != nil {
		logger.Error(ctx, "Error copying permission model to entity", err.Error())
		return nil, err
	}

	return &entityPerm, nil
}
