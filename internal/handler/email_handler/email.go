package email_handler

import "wtm-backend/internal/domain"

type EmailHandler struct {
	emailUsecase domain.EmailUsecase
}

func NewEmailHandler(emailUsecase domain.EmailUsecase) *EmailHandler {
	return &EmailHandler{
		emailUsecase: emailUsecase,
	}
}
