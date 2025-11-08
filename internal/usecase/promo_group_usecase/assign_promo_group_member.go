package promo_group_usecase

import (
	"context"
	"errors"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/promogroupdto"
	"wtm-backend/pkg/logger"
)

func (pgu *PromoGroupUsecase) AssignPromoGroupMember(ctx context.Context, req *promogroupdto.AssignPromoGroupMemberRequest) error {

	var dataMembers []entity.User

	//Check if the promo group exists
	promoGroup, err := pgu.promoGroupRepo.GetPromoGroupByID(ctx, req.PromoGroupID)
	if err != nil {
		logger.Error(ctx, "Error getting promo group", err.Error())
		return err
	}

	if promoGroup == nil {
		logger.Error(ctx, "Promo group not found", "promoGroupID", req.PromoGroupID)
		return errors.New("promo group not found")
	}

	if req.AgentCompanyID > 0 {

		members, _, err := pgu.userRepo.GetUsersByAgentCompany(ctx, req.AgentCompanyID, "", 0, 0)
		if err != nil {
			logger.Error(ctx, "Error getting members by agent company", err.Error())
			return err
		}

		if len(members) == 0 {
			logger.Error(ctx, "No members found for agent company", "agentCompanyID", req.AgentCompanyID)
			return errors.New("no members found for agent company")
		}

		dataMembers = members // Assuming you want to add the first member found for the agent company

	} else if req.MemberID > 0 {

		member, err := pgu.userRepo.GetUserByID(ctx, req.MemberID)
		if err != nil {
			logger.Error(ctx, "Error getting member", err.Error())
			return err
		}

		if member == nil {
			logger.Error(ctx, "Member not found", "memberID", req.MemberID)
			return errors.New("member not found")
		}

		dataMembers = append(dataMembers, *member) // Wrap the single member in a slice
	}

	var userIDsToUpdate []uint
	for _, member := range dataMembers {
		if member.PromoGroupID == nil || *member.PromoGroupID != req.PromoGroupID {
			userIDsToUpdate = append(userIDsToUpdate, member.ID)
		}
	}

	if len(userIDsToUpdate) == 0 {
		return errors.New("no members to update in the promo group")
	}

	err = pgu.userRepo.BulkUpdatePromoGroupMember(ctx, userIDsToUpdate, req.PromoGroupID)
	if err != nil {
		logger.Error(ctx, "Error updating promo group members", err.Error())
		return err
	}

	return nil
}
