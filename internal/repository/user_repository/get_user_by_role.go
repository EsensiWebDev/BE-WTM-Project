package user_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetUserByRole(ctx context.Context, roleID uint, search string, limit, page int) ([]entity.User, int64, error) {
	db := ur.db.GetTx(ctx)

	var users []model.User
	var total int64

	query := db.WithContext(ctx).
		Model(&model.User{}).
		Where("role_id = ?", roleID)

	if roleID == constant.DefaultRoleAgent {
		query = query.
			Preload("AgentCompany").
			Preload("PromoGroup").
			Preload("Status").
			Select("id, full_name, agent_company_id, promo_group_id, email, phone, kakao_talk_id, status_id")
	} else {
		query = query.
			Preload("Status").
			Select("id, full_name, email, phone, status_id")
	}

	if strings.TrimSpace(search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(search)
		query = query.Where("full_name ILIKE ?  ", "%"+safeSearch+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting users by role", err.Error())
		return nil, total, err
	}

	if limit > 0 {
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&users).Error; err != nil {
		logger.Error(ctx, "Error to get user by Id", err.Error())
		return nil, total, err
	}

	var entityUsers []entity.User
	for _, user := range users {
		var entityUser entity.User
		if err := utils.CopyPatch(&entityUser, user); err != nil {
			logger.Error(ctx, "Error copying user model to entity", err.Error())
			return nil, total, err
		}

		if user.AgentCompany != nil {
			entityUser.AgentCompanyName = user.AgentCompany.Name
		}

		if user.PromoGroup != nil {
			entityUser.PromoGroupName = user.PromoGroup.Name
		}

		if strings.TrimSpace(user.Status.Status) != "" {
			entityUser.StatusName = user.Status.Status
		}

		entityUsers = append(entityUsers, entityUser)
	}

	return entityUsers, total, nil
}
