package notifdto

import (
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto"
)

type ListNotificationsRequest struct {
	dto.PaginationRequest `json:",inline"`
}

type ListNotificationsResponse struct {
	Notifications []entity.Notification `json:"notifications"`
	Total         int64                 `json:"total"`
}
