package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Global logger instance
	log *zap.Logger
	// Ensure initialization happens only once
	once sync.Once
)

// Config defines logger configuration
type Config struct {
	Level      string
	OutputPath string
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int // days
	Compress   bool
	JSON       bool
}

// RotationConfig defines configuration for log rotation
type RotationConfig struct {
	Filename   string
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int // days
	Compress   bool
}

// DefaultConfig provides a default configuration for development
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		OutputPath: "",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
		JSON:       false,
	}
}

// DefaultRotationConfig provides a default configuration for log rotation
func DefaultRotationConfig() RotationConfig {
	return RotationConfig{
		Filename:   "blockchain-client.log",
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
}

// Init initializes the logger with the given configuration
func Init(cfg Config) *zap.Logger {
	once.Do(func() {
		// Setup output
		var sink zapcore.WriteSyncer
		if cfg.OutputPath == "" {
			sink = zapcore.AddSync(os.Stdout)
		} else {
			sink = zapcore.AddSync(&lumberjack.Logger{
				Filename:   cfg.OutputPath,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			})
		}

		// Set encoder
		var encoder zapcore.Encoder
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		if cfg.JSON {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		// Set level
		level := zap.InfoLevel
		switch cfg.Level {
		case "debug":
			level = zap.DebugLevel
		case "info":
			level = zap.InfoLevel
		case "warn":
			level = zap.WarnLevel
		case "error":
			level = zap.ErrorLevel
		}

		core := zapcore.NewCore(encoder, sink, zap.NewAtomicLevelAt(level))
		log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})

	return log
}

// InitWithRotation initializes the logger with file rotation
func InitWithRotation(level string, rotationCfg RotationConfig) *zap.Logger {
	once.Do(func() {
		// Configure rotating logger
		rotatingLogger := &lumberjack.Logger{
			Filename:   rotationCfg.Filename,
			MaxSize:    rotationCfg.MaxSize,
			MaxBackups: rotationCfg.MaxBackups,
			MaxAge:     rotationCfg.MaxAge,
			Compress:   rotationCfg.Compress,
		}

		// Configure writers - use zapcore.AddSync to properly wrap writers
		consoleSink := zapcore.AddSync(os.Stdout)
		fileSink := zapcore.AddSync(rotatingLogger)

		// Configure encoder
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		// Set level
		zapLevel := zap.InfoLevel
		switch level {
		case "debug":
			zapLevel = zap.DebugLevel
		case "info":
			zapLevel = zap.InfoLevel
		case "warn":
			zapLevel = zap.WarnLevel
		case "error":
			zapLevel = zap.ErrorLevel
		}
		atomicLevel := zap.NewAtomicLevelAt(zapLevel)

		// Create core for both console and file output
		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleSink, atomicLevel),
			zapcore.NewCore(fileEncoder, fileSink, atomicLevel),
		)

		// Create logger
		log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	})

	return log
}

// GetLogger returns the global logger instance, initializing with defaults if necessary
func GetLogger() *zap.Logger {
	if log == nil {
		return Init(DefaultConfig())
	}
	return log
}

// Sync flushes any buffered log entries
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

// Debug logs a debug message
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info logs an info message
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn logs a warning message
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error logs an error message
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// With returns a logger with additional fields
func With(fields ...zap.Field) *zap.Logger {
	return GetLogger().With(fields...)
}
