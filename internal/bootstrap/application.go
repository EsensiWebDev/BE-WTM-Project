package bootstrap

import (
	"context"
	"wtm-backend/config"
	"wtm-backend/internal/infrastructure/cache"
	"wtm-backend/internal/infrastructure/database"
	"wtm-backend/internal/infrastructure/email"
	"wtm-backend/internal/infrastructure/storage"
	"wtm-backend/internal/middleware"
	"wtm-backend/internal/usecase/auth_usecase"
	"wtm-backend/internal/usecase/banner_usecase"
	"wtm-backend/internal/usecase/booking_usecase"
	"wtm-backend/internal/usecase/email_usecase"
	"wtm-backend/internal/usecase/file_usecase"
	"wtm-backend/internal/usecase/hotel_usecase"
	"wtm-backend/internal/usecase/notification_usecase"
	"wtm-backend/internal/usecase/promo_group_usecase"
	"wtm-backend/internal/usecase/promo_usecase"
	"wtm-backend/internal/usecase/report_usecase"
	"wtm-backend/internal/usecase/user_usecase"
	"wtm-backend/pkg/logger"
)

type Application struct {
	Config        *config.Config
	Middleware    *middleware.Middleware
	Usecases      AppUsecases
	dB            *database.DBPostgre
	redis         *cache.RedisClient
	storageClient *storage.MultiStorageClient
	email         *email.SMTPEmailSender
}

type AppUsecases struct {
	AuthUsecase         *auth_usecase.AuthUsecase
	UserUsecase         *user_usecase.UserUsecase
	PromoUsecase        *promo_usecase.PromoUsecase
	HotelUsecase        *hotel_usecase.HotelUsecase
	BannerUsecase       *banner_usecase.BannerUsecase
	PromoGroupUsecase   *promo_group_usecase.PromoGroupUsecase
	BookingUsecase      *booking_usecase.BookingUsecase
	ReportUsecase       *report_usecase.ReportUsecase
	NotificationUsecase *notification_usecase.NotificationUsecase
	EmailUsecase        *email_usecase.EmailUsecase
	FileUsecase         *file_usecase.FileUsecase
}

func NewApplication() *Application {
	ctx := context.Background()
	logger.Info(ctx, "Initializing application...")

	// Initialize dependencies
	deps, err := initializeDependencies(ctx)
	if err != nil {
		logger.Fatal(ctx, "Dependency initialization failed", err.Error())
	}

	// Initialize repositories
	repos := initializeRepositories(deps)

	// Initialize usecases
	usecases := initializeUsecases(deps, repos)

	return &Application{
		Config:        deps.Config,
		Usecases:      usecases,
		Middleware:    deps.Middleware,
		dB:            deps.DB,
		redis:         deps.Redis,
		storageClient: deps.Storage,
		email:         deps.EmailSender,
	}
}
