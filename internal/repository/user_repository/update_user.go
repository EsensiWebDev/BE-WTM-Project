package user_repository

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	db := ur.db.GetTx(ctx)

	var modelUser model.User
	if err := utils.CopyStrict(&modelUser, user); err != nil {
		logger.Error(ctx, "Error copying user entity to model", err.Error())
		return nil, err
	}

	updateData := map[string]interface{}{
		"password":         modelUser.Password,
		"full_name":        modelUser.FullName,
		"username":         modelUser.Username,
		"email":            modelUser.Email,
		"phone":            modelUser.Phone,
		"kakao_talk_id":    modelUser.KakaoTalkID,
		"agent_company_id": modelUser.AgentCompanyID,
		"certificate":      modelUser.Certificate,
		"photo_selfie":     modelUser.PhotoSelfie,
		"photo_id_card":    modelUser.PhotoIDCard,
		"name_card":        modelUser.NameCard,
		"status_id":        modelUser.StatusID,
		"promo_group_id":   modelUser.PromoGroupID,
	}

	err := db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(updateData).Error
	if err != nil {
		logger.Error(ctx, "Error to update user", err.Error())
		return nil, err
	}

	return user, nil
}
