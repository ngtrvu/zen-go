package log

import (
	"io"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var logLevelSeverity = map[zerolog.Level]string{
	zerolog.DebugLevel: "DEBUG",
	zerolog.InfoLevel:  "INFO",
	zerolog.WarnLevel:  "WARNING",
	zerolog.ErrorLevel: "ERROR",
	zerolog.PanicLevel: "CRITICAL",
	zerolog.FatalLevel: "CRITICAL",
}

func InitZeroLog(logLevel string, writers ...io.Writer) {
	level, _ := zerolog.ParseLevel(logLevel)

	zerolog.SetGlobalLevel(level)
	zerolog.LevelFieldName = "severity"
	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return logLevelSeverity[l]
	}
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.CallerSkipFrameCount = 3
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()

	// setup cloud logging for log streaming
	multiWriters := zerolog.MultiLevelWriter(writers...)
	GLogger = zerolog.New(multiWriters).With().Timestamp().Logger()

	Info("zerolog is initialized, log level: %s", logLevel)
}

func Debug(message string, v ...interface{}) {
	log.Debug().Msgf(message, v...)
}

func Info(message string, v ...interface{}) {
	log.Info().Msgf(message, v...)
}

func Warn(message string, v ...interface{}) {
	log.Warn().Msgf(message, v...)
}

func Error(message string, v ...interface{}) {
	log.Error().Msgf(message, v...)
}

func Fatal(message string, v ...interface{}) {
	log.Fatal().Msgf(message, v...)
}

func Streaming(message string, v ...interface{}) {
	GLogger.Info().Msgf(message, v...)
}
