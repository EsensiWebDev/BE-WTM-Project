package utils

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"os"
	"strconv"
	"strings"
	"time"
	"wtm-backend/pkg/logger"
)

var appEnv = os.Getenv("APP_ENV") // "DEV", "STAGING", "PROD"

// getFromSSM fetches a parameter from AWS SSM with short timeout
func getFromSSM(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-southeast-1"))
	if err != nil {
		return "", err
	}

	client := ssm.NewFromConfig(cfg)
	param, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return *param.Parameter.Value, nil
}

// getEnvValue resolves value based on environment
func getEnvValue(key string) (string, error) {
	ctx := context.Background()

	if strings.ToUpper(appEnv) == "PROD" {
		// Only use SSM in production
		val, err := getFromSSM(key)
		if err != nil || val == "" {
			logger.Error(ctx, "Missing required SSM parameter:", key)
			return "", errors.New("missing required SSM parameter: " + key)
		}
		return val, nil
	}

	// DEV or STAGING: fallback to ENV
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val, nil
	}

	logger.Warn(ctx, "ENV variable not found:", key)
	return "", errors.New("missing ENV variable: " + key)
}

// GetStringEnv returns string or default
func GetStringEnv(key string, defaultVal string) string {
	val, err := getEnvValue(key)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetIntEnv returns int or default
func GetIntEnv(key string, defaultVal int) int {
	val, err := getEnvValue(key)
	if err != nil {
		return defaultVal
	}
	if parsed, err := strconv.Atoi(val); err == nil {
		return parsed
	}
	return defaultVal
}

// GetBoolEnv returns bool or default
func GetBoolEnv(key string, defaultVal bool) bool {
	val, err := getEnvValue(key)
	if err != nil {
		return defaultVal
	}
	if parsed, err := strconv.ParseBool(val); err == nil {
		return parsed
	}
	return defaultVal
}

// GetDurationEnv returns duration or default
func GetDurationEnv(key string, defaultVal time.Duration) time.Duration {
	val, err := getEnvValue(key)
	if err != nil {
		return defaultVal
	}
	if parsed, err := parseFlexibleDuration(val); err == nil {
		return parsed
	}
	return defaultVal
}

// parseFlexibleDuration handles "s", "m", "h", "d"
func parseFlexibleDuration(input string) (time.Duration, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	if strings.HasSuffix(input, "d") {
		numStr := strings.TrimSuffix(input, "d")
		if num, err := strconv.Atoi(numStr); err == nil {
			return time.Duration(num) * 24 * time.Hour, nil
		}
	}
	return time.ParseDuration(input)
}
