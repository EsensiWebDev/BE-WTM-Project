package promo_usecase

import (
	"context"
	"fmt"
	"wtm-backend/internal/dto/promodto"
	"wtm-backend/internal/repository/filter"
	"wtm-backend/pkg/logger"
)

func (pu *PromoUsecase) ListPromosForAgent(ctx context.Context, req *promodto.ListPromosForAgentRequest) (*promodto.ListPromosForAgentResponse, error) {
	// Get agent Id from context
	userCtx, err := pu.middleware.GenerateUserFromContext(ctx)
	if err != nil {
		logger.Error(ctx, "failed to get user from context", err.Error())
		return nil, err
	}

	if userCtx == nil {
		logger.Error(ctx, "user context is nil")
		return nil, err
	}

	agentID := userCtx.ID

	filterReq := &filter.PromoFilter{
		PaginationRequest: req.PaginationRequest,
		AgentID:           agentID,
	}

	promos, total, err := pu.promoRepo.GetPromosWithHotels(ctx, filterReq)
	if err != nil {
		logger.Error(ctx, "failed to get promos", err.Error())
		return nil, err
	}

	dataPromos := make([]promodto.PromosForAgent, 0, len(promos))
	for _, promo := range promos {
		var dataHotels []string
		for _, roomType := range promo.PromoRoomTypes {
			dataHotels = append(dataHotels, fmt.Sprintf("%s %s", roomType.HotelName, roomType.Province))
		}
		dataPromos = append(dataPromos, promodto.PromosForAgent{
			ID:          promo.ID,
			Name:        promo.Name,
			Code:        promo.Code,
			Description: promo.Description,
			Hotel:       dataHotels,
		})
	}

	return &promodto.ListPromosForAgentResponse{
		Data:  dataPromos,
		Total: total,
	}, nil
}
