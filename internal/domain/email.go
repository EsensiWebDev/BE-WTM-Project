package domain

import (
	"context"
	"wtm-backend/internal/domain/entity"
	"wtm-backend/internal/dto/emaildto"
	"wtm-backend/internal/repository/filter"
)

type EmailSender interface {
	Send(ctx context.Context, to, subject, bodyHTML, bodyText string) error
}

type EmailRepository interface {
	GetEmailTemplateByName(ctx context.Context, name string) (*entity.EmailTemplate, error)
	UpdateEmailTemplate(ctx context.Context, template *entity.EmailTemplate) error
	CreateEmailLog(ctx context.Context, log *entity.EmailLog) error
	GetEmailLogs(ctx context.Context, filter filter.DefaultFilter) ([]entity.EmailLog, int64, error)
}

type EmailUsecase interface {
	EmailTemplate(ctx context.Context) (*emaildto.EmailTemplateResponse, error)
	UpdateEmailTemplate(ctx context.Context, req *emaildto.UpdateEmailTemplateRequest) error
	SendContactUsEmail(ctx context.Context, req *emaildto.SendContactUsEmailRequest) error
	ListEmailLogs(ctx context.Context, req *emaildto.ListEmailLogsRequest) (*emaildto.ListEmailLogsResponse, error)
}
