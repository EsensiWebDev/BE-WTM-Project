package user_repository

import (
	"context"
	"strings"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/infrastructure/database/model"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
	"wtm-backend/pkg/utils"
)

func (ur *UserRepository) GetUsers(ctx context.Context, filter filter.UserFilter) ([]entity.User, int64, error) {
	db := ur.db.GetTx(ctx)
	query := db.WithContext(ctx).Model(&model.User{}).Debug()

	// Select default fields
	selectFields := []string{"id", "full_name", "email", "phone", "status_id"}

	// Apply filters
	if filter.RoleID != nil {
		query = query.Where("role_id = ?", *filter.RoleID).Preload("Status")

		if *filter.RoleID == constant.DefaultRoleAgent {
			selectFields = append(selectFields, "agent_company_id")
			query = query.Preload("AgentCompany")

			if filter.Scope == constant.ScopeControl {
				selectFields = append(selectFields, "agent_company_id", "kakao_talk_id", "photo_selfie", "certificate", "name_card", "photo_id_card")
			} else if filter.Scope == constant.ScopeManagement {
				selectFields = append(selectFields, "agent_company_id", "promo_group_id")
				query = query.Where("status_id = ?", constant.StatusUserActiveID).Preload("PromoGroup")
			}

		}
	}

	if filter.StatusID != nil {
		if *filter.StatusID > 0 {
			query = query.Where("status_id = ?", *filter.StatusID)
		}
	}

	if filter.AgentCompanyID != nil {
		query = query.Where("agent_company_id = ?", *filter.AgentCompanyID)
	}

	query = query.Select(selectFields)

	// Search
	if strings.TrimSpace(filter.Search) != "" {
		safeSearch := utils.EscapeAndNormalizeSearch(filter.Search)
		query = query.Where("LOWER(full_name) ILIKE ? ESCAPE '\\' ", "%"+safeSearch+"%")
	}

	// Count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error(ctx, "Error counting users", err.Error())
		return nil, total, err
	}

	// Pagination
	if filter.Limit > 0 {
		if filter.Page < 1 {
			filter.Page = 1
		}
		offset := (filter.Page - 1) * filter.Limit
		query = query.Limit(filter.Limit).Offset(offset)
	}

	// Execute
	var users []model.User
	if err := query.Find(&users).Error; err != nil {
		logger.Error(ctx, "Error fetching users", err.Error())
		return nil, total, err
	}

	// Mapping
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
