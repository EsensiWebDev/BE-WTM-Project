package config

import (
	"github.com/joho/godotenv"
	"strings"
	"time"
	"wtm-backend/pkg/utils"
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

	URL string

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

	EmailHost     string
	EmailPort     int
	EmailUser     string
	EmailPass     string
	EmailFrom     string
	EmailProvider string

	SupportEmail string

	AWSConfig AWSConfig
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

		URL: utils.GetStringEnv("URL", ""),

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

		EmailHost:     utils.GetStringEnv("EMAIL_HOST", "localhost"),
		EmailPort:     utils.GetIntEnv("EMAIL_PORT", 1025),
		EmailUser:     utils.GetStringEnv("EMAIL_USERNAME", ""),
		EmailPass:     utils.GetStringEnv("EMAIL_PASSWORD", ""),
		EmailFrom:     utils.GetStringEnv("EMAIL_FROM", "noreply@dev.local"),
		EmailProvider: utils.GetStringEnv("EMAIL_PROVIDER", "mailhog"),

		SupportEmail: utils.GetStringEnv("SUPPORT_EMAIL", "support@dev.local"),

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
