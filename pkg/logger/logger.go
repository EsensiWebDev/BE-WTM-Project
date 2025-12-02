package logger

import (
	"context"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = logrus.New()

func InitLogger() {

	// Set format log dalam JSON
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	// Lumberjack log rotation config
	rotatingLogger := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10,   // MB
		MaxBackups: 5,    // backup files
		MaxAge:     28,   // days
		Compress:   true, // gzip compressed
	}

	// Combine stdout + file
	multiWriter := io.MultiWriter(os.Stdout, rotatingLogger)
	log.SetOutput(multiWriter)

	// Atur level log sesuai kebutuhan (Info, Warn, Error, dsb.)
	log.SetLevel(logrus.InfoLevel)
}

// Info logs informational messages.
func Info(ctx context.Context, message string, args ...interface{}) {
	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "unknown"
	}
	log.WithFields(logrus.Fields{
		"traceID": traceID,
		"func":    getCallerFuncName(),
		"data":    args,
	}).Info(message)
}

// Warn logs warning messages.
func Warn(ctx context.Context, message string, args ...interface{}) {
	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "unknown"
	}
	log.WithFields(logrus.Fields{
		"traceID": traceID,
		"func":    getCallerFuncName(),
		"data":    args,
	}).Warn(message)
}

// Error logs error messages.
func Error(ctx context.Context, message string, args ...interface{}) {
	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "unknown"
	}
	log.WithFields(logrus.Fields{
		"traceID": traceID,
		"func":    getCallerFuncName(),
		"data":    args,
	}).Error(message)
}

// Fatal logs fatal errors and exits.
func Fatal(ctx context.Context, message string, args ...interface{}) {
	traceID := ctx.Value("traceID")
	if traceID == nil {
		traceID = "unknown"
	}
	log.WithFields(logrus.Fields{
		"func": getCallerFuncName(),
		"data": args,
	}).Fatal(message)
}

func getCallerFuncName() string {
	pc, _, _, ok := runtime.Caller(2) // 2 tingkat di atas, agar dapet pemanggil log
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	// Ambil hanya nama fungsi terakhir (tanpa package path)
	parts := strings.Split(fn.Name(), ".")
	return parts[len(parts)-1]
}
