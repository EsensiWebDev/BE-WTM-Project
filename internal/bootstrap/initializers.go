package bootstrap

import (
	"context"
	"fmt"
	"wtm-backend/config"
	"wtm-backend/internal/infrastructure/cache"
	"wtm-backend/internal/infrastructure/database"
	"wtm-backend/internal/infrastructure/email"
	"wtm-backend/internal/infrastructure/storage"
	"wtm-backend/internal/middleware"
	"wtm-backend/internal/repository/auth_repository"
	"wtm-backend/internal/repository/banner_repository"
	"wtm-backend/internal/repository/booking_repository"
	"wtm-backend/internal/repository/driver"
	"wtm-backend/internal/repository/email_repository"
	"wtm-backend/internal/repository/hotel_repository"
	"wtm-backend/internal/repository/notification_repository"
	"wtm-backend/internal/repository/promo_group_repository"
	"wtm-backend/internal/repository/promo_repository"
	"wtm-backend/internal/repository/report_repository"
	"wtm-backend/internal/repository/user_repository"
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

func initializeDependencies(ctx context.Context) (*Dependencies, error) {
	cfg := config.LoadConfig()

	db, err := database.NewDBPostgre(cfg)
	if err != nil {
		logger.Error(ctx, "failed to initialize database", err.Error())
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		logger.Error(ctx, "failed to initialize Redis client", err.Error())
		return nil, fmt.Errorf("failed to initialize Redis client: %w", err)
	}

	storageClient, err := storage.NewMultiStorageClient(cfg)
	if err != nil {
		logger.Error(ctx, "failed to initialize storage client", err.Error())
		return nil, fmt.Errorf("failed to initialize storage client: %w", err)
	}

	emailClient := email.NewSMTPEmailSender(cfg)

	if storageClient.ActiveStorage == nil {
		logger.Error(ctx, "no active storage client available")
		return nil, fmt.Errorf("no active storage client available")
	}

	newMiddleware := middleware.NewMiddleware(cfg, auth_repository.NewAuthRepository(redisClient, db))
	dbTransaction := driver.NewDatabaseTransaction(db)

	return &Dependencies{
		Config:        cfg,
		DB:            db,
		Redis:         redisClient,
		Storage:       storageClient,
		EmailSender:   emailClient,
		Middleware:    newMiddleware,
		DBTransaction: dbTransaction,
	}, nil
}

func initializeRepositories(deps *Dependencies) *Repositories {
	return &Repositories{
		UserRepo:         user_repository.NewUserRepository(deps.DB),
		AuthRepo:         auth_repository.NewAuthRepository(deps.Redis, deps.DB),
		PromoRepo:        promo_repository.NewPromoRepository(deps.DB),
		PromoGroupRepo:   promo_group_repository.NewPromoGroupRepository(deps.DB),
		HotelRepo:        hotel_repository.NewHotelRepository(deps.DB),
		BannerRepo:       banner_repository.NewBannerRepository(deps.DB),
		BookingRepo:      booking_repository.NewBookingRepository(deps.DB, deps.Redis),
		EmailRepo:        email_repository.NewEmailRepository(deps.DB),
		ReportRepo:       report_repository.NewReportRepository(deps.DB),
		NotificationRepo: notification_repository.NewNotificationRepository(deps.DB),
	}
}

func initializeUsecases(deps *Dependencies, repos *Repositories) AppUsecases {
	storageActive := deps.Storage.ActiveStorage

	return AppUsecases{
		AuthUsecase:         auth_usecase.NewAuthUsecase(repos.UserRepo, repos.AuthRepo, deps.Config, storageActive, deps.Middleware, deps.EmailSender, repos.EmailRepo, deps.DBTransaction),
		UserUsecase:         user_usecase.NewUserUsecase(repos.UserRepo, repos.AuthRepo, repos.PromoGroupRepo, repos.EmailRepo, deps.Config, storageActive, deps.Middleware, deps.DBTransaction, deps.EmailSender),
		PromoUsecase:        promo_usecase.NewPromoUsecase(repos.PromoRepo, deps.DBTransaction, deps.Middleware),
		HotelUsecase:        hotel_usecase.NewHotelUsecase(repos.HotelRepo, storageActive, deps.DBTransaction, deps.Config),
		BannerUsecase:       banner_usecase.NewBannerUsecase(repos.BannerRepo, deps.DBTransaction, storageActive),
		PromoGroupUsecase:   promo_group_usecase.NewPromoGroupUsecase(repos.PromoGroupRepo, repos.UserRepo),
		BookingUsecase:      booking_usecase.NewBookingUsecase(repos.BookingRepo, repos.HotelRepo, repos.PromoRepo, deps.Middleware, deps.DBTransaction, storageActive, deps.Config, repos.EmailRepo, deps.EmailSender, repos.UserRepo),
		ReportUsecase:       report_usecase.NewReportUsecase(repos.ReportRepo),
		NotificationUsecase: notification_usecase.NewNotificationUsecase(repos.NotificationRepo, deps.Middleware, deps.DBTransaction),
		EmailUsecase:        email_usecase.NewEmailUsecase(repos.EmailRepo, deps.EmailSender, repos.BookingRepo),
		FileUsecase:         file_usecase.NewFileUsecase(storageActive),
	}
}
