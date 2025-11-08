package bootstrap

import (
	"context"
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
	logger.InitLogger()
	ctx := context.Background()
	logger.Info(ctx, "Initializing application...")

	// Load config
	cfg := config.LoadConfig()

	// Init database
	db, err := database.NewDBPostgre(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database", err.Error())
	}

	// Initialize Redis client
	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize Redis client", err.Error())
	}

	// Initialize storage client
	storageClient, err := storage.NewMultiStorageClient(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize storage client", err.Error())
	}

	// Init Email
	emailClient := email.NewSMTPEmailSender(cfg)

	// Inisialisasi Repository
	userRepo := user_repository.NewUserRepository(db)
	authRepo := auth_repository.NewAuthRepository(redisClient, db)
	promoRepo := promo_repository.NewPromoRepository(db)
	promoGroupRepo := promo_group_repository.NewPromoGroupRepository(db)
	hotelRepo := hotel_repository.NewHotelRepository(db)
	dbTransaction := driver.NewDatabaseTransaction(db)
	bannerRepo := banner_repository.NewBannerRepository(db)
	bookingRepo := booking_repository.NewBookingRepository(db, redisClient)
	emailRepo := email_repository.NewEmailRepository(db)
	reportRepo := report_repository.NewReportRepository(db)
	notifRepo := notification_repository.NewNotificationRepository(db)

	// Initialize Middleware
	newMiddleware := middleware.NewMiddleware(cfg, authRepo)

	emailSender := email.NewSMTPEmailSender(cfg)

	storageClientActive := storageClient.ActiveStorage
	if storageClientActive == nil {
		logger.Fatal("No active storage client available")
	}

	// Inisialisasi UseCase
	authUsecase := auth_usecase.NewAuthUsecase(userRepo, authRepo, cfg, storageClientActive, newMiddleware, emailSender, emailRepo, dbTransaction)
	userUsecase := user_usecase.NewUserUsecase(userRepo, authRepo, promoGroupRepo, emailRepo, cfg, storageClientActive, newMiddleware, dbTransaction, emailSender)
	promoUsecase := promo_usecase.NewPromoUsecase(promoRepo, dbTransaction)
	hotelUsecase := hotel_usecase.NewHotelUsecase(hotelRepo, storageClientActive, dbTransaction, cfg)
	bannerUsecase := banner_usecase.NewBannerUsecase(bannerRepo, dbTransaction, storageClientActive)
	promoGroupUsecase := promo_group_usecase.NewPromoGroupUsecase(promoGroupRepo, userRepo)
	bookingUsecase := booking_usecase.NewBookingUsecase(bookingRepo, hotelRepo, promoRepo, newMiddleware, dbTransaction, storageClientActive, cfg, emailRepo, emailSender)
	reportUsecase := report_usecase.NewReportUsecase(reportRepo)
	notificationUsecase := notification_usecase.NewNotificationUsecase(notifRepo, newMiddleware, dbTransaction)
	emailUsecase := email_usecase.NewEmailUsecase(emailRepo, emailSender, bookingRepo)
	fielUsecase := file_usecase.NewFileUsecase(storageClientActive)

	// Register Usecases and Repositories
	appUsecases := AppUsecases{
		AuthUsecase:         authUsecase,
		UserUsecase:         userUsecase,
		PromoUsecase:        promoUsecase,
		HotelUsecase:        hotelUsecase,
		BannerUsecase:       bannerUsecase,
		PromoGroupUsecase:   promoGroupUsecase,
		BookingUsecase:      bookingUsecase,
		ReportUsecase:       reportUsecase,
		NotificationUsecase: notificationUsecase,
		EmailUsecase:        emailUsecase,
		FileUsecase:         fielUsecase,
	}

	return &Application{
		Config:        cfg,
		Usecases:      appUsecases,
		Middleware:    newMiddleware,
		dB:            db,
		redis:         redisClient,
		storageClient: storageClient,
		email:         emailClient,
	}
}
