package bootstrap

import (
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
)

type Dependencies struct {
	Config        *config.Config
	DB            *database.DBPostgre
	Redis         *cache.RedisClient
	Storage       *storage.MultiStorageClient
	EmailSender   *email.SMTPEmailSender
	Middleware    *middleware.Middleware
	DBTransaction *driver.DatabaseTransaction
}

type Repositories struct {
	UserRepo         *user_repository.UserRepository
	AuthRepo         *auth_repository.AuthRepository
	PromoRepo        *promo_repository.PromoRepository
	PromoGroupRepo   *promo_group_repository.PromoGroupRepository
	HotelRepo        *hotel_repository.HotelRepository
	BannerRepo       *banner_repository.BannerRepository
	BookingRepo      *booking_repository.BookingRepository
	EmailRepo        *email_repository.EmailRepository
	ReportRepo       *report_repository.ReportRepository
	NotificationRepo *notification_repository.NotificationRepository
}
