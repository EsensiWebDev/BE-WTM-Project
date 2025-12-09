package config

import (
	"strings"
	"time"
	"wtm-backend/pkg/utils"

	"github.com/joho/godotenv"
)

type AWSConfig struct {
	Region            string
	AccessKeyID       string
	SecretAccessKey   string
	Endpoint          string // untuk local development/test
	S3BucketURL       string
	S3PresignDuration time.Duration
}

type Config struct {
	AppEnv     string
	Host       string
	ServerPort string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresName     string

	JWTSecret     string
	RefreshSecret string

	URL        string
	URLFEAgent string
	URLFEAdmin string

	HostRedis     string
	PortRedis     string
	PasswordRedis string

	MinioHost         string
	MinioPort         string
	MinioRootUser     string
	MinioRootPassword string

	SecureService bool

	AllowOrigins   string
	AllowedOrigins []string

	DurationCtxTOFast       time.Duration
	DurationCtxTOSlow       time.Duration
	DurationCtxTOFile       time.Duration
	DurationAccessFileMinio time.Duration
	DurationAccessToken     time.Duration
	DurationRefreshToken    time.Duration
	DurationMaxAgeCORS      time.Duration
	DurationLinkExpiration  time.Duration

	DefaultCancellationPeriod int
	DefaultCheckInHour        string
	DefaultCheckOutHour       string

	HostSES               string
	UsernameSES           string
	PortSES               int
	PasswordSES           string
	DisableAuthSES        bool
	UseTLSSES             bool
	HostNameSES           string
	DefaultFromSES        string
	SupportEmailSES       string
	ProviderSESReplyTo    bool
	ProviderSESReturnPath bool
	TimeoutSES            time.Duration

	HostGmail               string
	UsernameGmail           string
	PortGmail               int
	PasswordGmail           string
	DisableAuthGmail        bool
	UseTLSGmail             bool
	HostNameGmail           string
	DefaultFromGmail        string
	SupportEmailGmail       string
	ProviderGmailReplyTo    bool
	ProviderGmailReturnPath bool
	TimeoutGmail            time.Duration

	HostMailhog        string
	PortMailhog        int
	DisableAuthMailhog bool
	UseTLSMailhog      bool
	DefaultFromMailhog string

	EmailContactUs string
	EmailFromAgent string
	EmailFromHotel string

	AWSConfig   AWSConfig
	AutoMigrate bool

	ProviderAgent string
	ProviderHotel string

	RetryDelay     time.Duration
	DialTimeout    time.Duration
	SendTimeout    time.Duration
	CommandTimeout time.Duration

	HostIP string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	config := &Config{
		AppEnv:     utils.GetStringEnv("APP_ENV", "DEV"),
		Host:       utils.GetStringEnv("HOST", "localhost:4816"),
		ServerPort: utils.GetStringEnv("SERVER_PORT", "4816"),

		PostgresHost:     utils.GetStringEnv("POSTGRES_HOST", ""),
		PostgresPort:     utils.GetStringEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     utils.GetStringEnv("POSTGRES_USER", "youruser"),
		PostgresPassword: utils.GetStringEnv("POSTGRES_PASSWORD", "yourpassword"),
		PostgresName:     utils.GetStringEnv("POSTGRES_DATABASE", "yourdatabase"),

		JWTSecret:     utils.GetStringEnv("JWT_SECRET", "yoursecretkey"),
		RefreshSecret: utils.GetStringEnv("REFRESH_SECRET", "yourrefreshkey"),

		URL:        utils.GetStringEnv("URL", ""),
		URLFEAgent: utils.GetStringEnv("URL_FE_AGENT", ""),
		URLFEAdmin: utils.GetStringEnv("URL_FE_ADMIN", ""),

		HostRedis:     utils.GetStringEnv("REDIS_HOST", "localhost"),
		PortRedis:     utils.GetStringEnv("REDIS_PORT", "6379"),
		PasswordRedis: utils.GetStringEnv("REDIS_PASSWORD", ""),

		MinioHost:         utils.GetStringEnv("MINIO_HOST", "localhost"),
		MinioPort:         utils.GetStringEnv("MINIO_PORT", "9000"),
		MinioRootUser:     utils.GetStringEnv("MINIO_ROOT_USER", "minioadmin"),
		MinioRootPassword: utils.GetStringEnv("MINIO_ROOT_PASSWORD", "minioadmin"),

		SecureService: utils.GetBoolEnv("SECURE_SERVICE", true),
		AllowOrigins:  utils.GetStringEnv("ALLOW_ORIGINS", "*"),

		DurationCtxTOFast:       utils.GetDurationEnv("CTX_TIMEOUT_FAST", 5*time.Second),
		DurationCtxTOSlow:       utils.GetDurationEnv("CTX_TIMEOUT_SLOW", 10*time.Second),
		DurationCtxTOFile:       utils.GetDurationEnv("CTX_TIMEOUT_FILE", 60*time.Second),
		DurationAccessFileMinio: utils.GetDurationEnv("EXPIRATION_PRIVATE_FILE", 5*time.Minute),
		DurationAccessToken:     utils.GetDurationEnv("EXPIRATION_TIME_ACCESS_TOKEN", 30*time.Minute),
		DurationRefreshToken:    utils.GetDurationEnv("EXPIRATION_TIME_REFRESH_TOKEN", 7*24*time.Hour),
		DurationMaxAgeCORS:      utils.GetDurationEnv("MAX_AGE_CORS", 12*time.Hour),
		DurationLinkExpiration:  utils.GetDurationEnv("DURATION_LINK_EXPIRATION", 45*time.Minute),

		DefaultCancellationPeriod: utils.GetIntEnv("DEFAULT_CANCEL_PERIOD", 5),
		DefaultCheckInHour:        utils.GetStringEnv("DEFAULT_CHECK_IN_HOUR", "14:00"),
		DefaultCheckOutHour:       utils.GetStringEnv("DEFAULT_CHECK_OUT_HOUR", "12:00"),

		HostSES:               utils.GetStringEnv("HOST_SES", ""),
		UsernameSES:           utils.GetStringEnv("USERNAME_SES", ""),
		PortSES:               utils.GetIntEnv("PORT_SES", 587),
		PasswordSES:           utils.GetStringEnv("PASSWORD_SES", ""),
		DisableAuthSES:        utils.GetBoolEnv("DISABLE_AUTH_SES", false),
		UseTLSSES:             utils.GetBoolEnv("USE_TLS_SES", false),
		HostNameSES:           utils.GetStringEnv("HOST_NAME_SES", ""),
		DefaultFromSES:        utils.GetStringEnv("DEFAULT_FROM_SES", ""),
		SupportEmailSES:       utils.GetStringEnv("SUPPORT_EMAIL_SES", ""),
		ProviderSESReplyTo:    utils.GetBoolEnv("PROVIDER_SES_REPLY_TO", false),
		ProviderSESReturnPath: utils.GetBoolEnv("PROVIDER_SES_RETURN_PATH", false),
		TimeoutSES:            utils.GetDurationEnv("TIMEOUT_SES", 12*time.Second),

		HostGmail:               utils.GetStringEnv("HOST_GMAIL", ""),
		UsernameGmail:           utils.GetStringEnv("USERNAME_GMAIL", ""),
		PortGmail:               utils.GetIntEnv("PORT_GMAIL", 587),
		PasswordGmail:           utils.GetStringEnv("PASSWORD_GMAIL", ""),
		DisableAuthGmail:        utils.GetBoolEnv("DISABLE_AUTH_GMAIL", false),
		UseTLSGmail:             utils.GetBoolEnv("USE_TLS_GMAIL", false),
		DefaultFromGmail:        utils.GetStringEnv("DEFAULT_FROM_GMAIL", ""),
		SupportEmailGmail:       utils.GetStringEnv("SUPPORT_EMAIL_GMAIL", ""),
		ProviderGmailReplyTo:    utils.GetBoolEnv("PROVIDER_GMAIL_REPLY_TO", false),
		ProviderGmailReturnPath: utils.GetBoolEnv("PROVIDER_GMAIL_RETURN_PATH", false),
		TimeoutGmail:            utils.GetDurationEnv("TIMEOUT_GMAIL", 12*time.Second),

		HostMailhog:        utils.GetStringEnv("HOST_MAILHOG", ""),
		PortMailhog:        utils.GetIntEnv("PORT_MAILHOG", 8025),
		DisableAuthMailhog: utils.GetBoolEnv("DISABLE_AUTH_MAILHOG", false),
		UseTLSMailhog:      utils.GetBoolEnv("USE_TLS_MAILHOG", false),
		DefaultFromMailhog: utils.GetStringEnv("DEFAULT_FROM_MAILHOG", ""),

		EmailContactUs: utils.GetStringEnv("EMAIL_CONTACT_US", "contact@wtm.com"),
		EmailFromAgent: utils.GetStringEnv("EMAIL_FROM_AGENT", "agent@wtm.com"),
		EmailFromHotel: utils.GetStringEnv("EMAIL_FROM_HOTEL", "admin@wtm.com"),

		ProviderAgent: utils.GetStringEnv("PROVIDER_AGENT", "mailhog"),
		ProviderHotel: utils.GetStringEnv("PROVIDER_HOTEL", "mailhog"),

		RetryDelay:     utils.GetDurationEnv("RETRY_DELAY", 2*time.Second),
		DialTimeout:    utils.GetDurationEnv("DIAL_TIMEOUT", 10*time.Second),
		SendTimeout:    utils.GetDurationEnv("SEND_TIMEOUT", 30*time.Second),
		CommandTimeout: utils.GetDurationEnv("COMMAND_TIMEOUT", 5*time.Second),

		HostIP: utils.GetStringEnv("HOST_IP", "127.0.0.1"),

		AutoMigrate: utils.GetBoolEnv("AUTO_MIGRATE", false),

		AWSConfig: AWSConfig{
			Region:            utils.GetStringEnv("AWS_REGION", "us-east-1"),
			AccessKeyID:       utils.GetStringEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey:   utils.GetStringEnv("AWS_SECRET_ACCESS_KEY", ""),
			Endpoint:          utils.GetStringEnv("AWS_ENDPOINT", ""), // untuk local development/test
			S3BucketURL:       utils.GetStringEnv("AWS_S3_BUCKET_URL", ""),
			S3PresignDuration: utils.GetDurationEnv("AWS_S3_PRESIGN_DURATION", 15*time.Minute),
		},
	}

	// Process allowed origins list
	config.AllowedOrigins = strings.Split(config.AllowOrigins, ",")

	// Optional override for dev
	if config.AppEnv == "DEV" {
		config.URL = ""
		config.SecureService = false
	}

	return config
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.AppEnv) == "production"
}
