package log

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type logger struct {
	SlowThreshold         time.Duration
	SourceField           string
	SkipErrRecordNotFound bool
	Logger                zerolog.Logger
}

func NewGormLogger(logLevel string) *logger {
	zeroLogger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// set zerologger level
	level, _ := zerolog.ParseLevel(logLevel)
	zeroLogger = zeroLogger.Level(level)

	return &logger{
		Logger:                zeroLogger,
		SkipErrRecordNotFound: true,
		SlowThreshold:         time.Millisecond * 500, // 0.5s consider as slow, should be warning
	}
}

func NewWithGormLogger(l zerolog.Logger, config gormlogger.Config) *logger {
	return &logger{
		Logger:                l,
		SkipErrRecordNotFound: config.IgnoreRecordNotFoundError,
		SlowThreshold:         config.SlowThreshold,
	}
}

func (l *logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Info().Msgf(s, args...)
}

func (l *logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Warn().Msgf(s, args...)
}

func (l *logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.Logger.Error().Msgf(s, args...)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := map[string]interface{}{
		"sql":      sql,
		"duration": elapsed,
	}

	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	// log warning if the error is "context canceled"
	if err != nil && errors.Is(err, context.Canceled) {
		l.Logger.Warn().Fields(fields).Msgf("[GORM] context canceled")
		return
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.SkipErrRecordNotFound) {
		l.Logger.Error().Err(err).Fields(fields).Msg("[GORM] query error")
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.Logger.Warn().Fields(fields).Msgf("[GORM] slow query")
		return
	}

	l.Logger.Debug().Fields(fields).Msgf("[GORM] query")
}
