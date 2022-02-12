package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/Yangiboev/simple-server-crytocurrency-api/config"
)

// For mapping config levels to application levels
var zapLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// Zap structure
type Zap struct {
	*zap.SugaredLogger
}

// New creates new zap instance
func New(mode string, cfg *config.Logger) *Zap {
	return &Zap{
		SugaredLogger: getSugar(mode, cfg),
	}
}

// Get sugar
func getSugar(mode string, cfg *config.Logger) *zap.SugaredLogger {
	var (
		level      = getLevel(cfg.Level)
		writer     = zapcore.AddSync(os.Stderr)
		encoderCfg = getEncoderConfig(mode)
		encoder    = getEncoder(cfg.Encoding, encoderCfg)

		core = zapcore.NewCore(encoder, writer, zap.NewAtomicLevelAt(level))

		options = make([]zap.Option, 0)
	)

	if !cfg.DisableCaller {
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(1))
	}

	if !cfg.DisableStacktrace {
		options = append(options, zap.AddStacktrace(level))
	}

	zapInstance := zap.New(core, options...)

	return zapInstance.Sugar()
}

// getEncoder returns encoder by application mode
func getEncoder(encoding string, encodingConfig zapcore.EncoderConfig) zapcore.Encoder {
	if encoding == config.LogConsoleEncoding {
		return zapcore.NewConsoleEncoder(encodingConfig)
	}
	return zapcore.NewJSONEncoder(encodingConfig)
}

// getEncoderConfig returns encoding configurations
func getEncoderConfig(mode string) zapcore.EncoderConfig {
	var encoderConfig zapcore.EncoderConfig

	if mode == config.DevelopmentMode {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}

	encoderConfig.LevelKey = "LEVEL"
	encoderConfig.CallerKey = "CALLER"
	encoderConfig.TimeKey = "TIME"
	encoderConfig.NameKey = "NAME"
	encoderConfig.MessageKey = "MESSAGE"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return encoderConfig
}

// getLevel returns logger level by application config level
func getLevel(level string) zapcore.Level {
	zapLevel, exist := zapLevelMap[level]
	if !exist {
		return zapcore.DebugLevel
	}

	return zapLevel
}
