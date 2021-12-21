package logger

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *Logger

const (
	timeFormat = time.RFC3339Nano
)

type LogLevel uint8

const (
	LogLevelDebug      LogLevel = 0
	LogLevelProduction LogLevel = 1
)

type Logger struct {
	serviceName string
	logger      *zap.SugaredLogger
	showHeaders bool
	showBody    bool
}

func NewLogger(serviceName string, logLevel LogLevel, encodeToJSON, showHeaders, showBody bool) *Logger {
	lvl := zapcore.DebugLevel
	if logLevel == LogLevelProduction {
		lvl = zapcore.InfoLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.UTC().Format(timeFormat))
	}

	var encoder zapcore.Encoder
	if encodeToJSON {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	if serviceName != "" {
		encoder.AddString("service_name", serviceName)
	}

	l := zap.New(zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(lvl),
	)).WithOptions(
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	)

	return &Logger{
		logger:      l.Sugar(),
		serviceName: serviceName,
		showHeaders: showHeaders,
		showBody:    showBody,
	}
}

func (l *Logger) Fatal(ctx context.Context, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Fatal(args...)
}

func (l *Logger) Fatalf(ctx context.Context, template string, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Fatalf(template, args...)
}

func (l *Logger) Error(ctx context.Context, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Error(args...)
}

func (l *Logger) Errorf(ctx context.Context, template string, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Errorf(template, args...)
}

func (l *Logger) Warn(ctx context.Context, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Warn(args...)
}

func (l *Logger) Warnf(ctx context.Context, template string, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Warnf(template, args...)
}

func (l *Logger) Info(ctx context.Context, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Info(args...)
}

func (l *Logger) Infof(ctx context.Context, template string, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Infof(template, args...)
}

func (l *Logger) Debug(ctx context.Context, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Debug(args...)
}

func (l *Logger) Debugf(ctx context.Context, template string, args ...interface{}) {
	appendHTTPInfo(ctx, l).logger.Debugf(template, args...)
}

func appendHTTPInfo(ctx context.Context, l *Logger) *Logger {
	if reqInfo := getRequestInfo(ctx); reqInfo != nil {
		// sort http headers
		headerNames := make([]string, 0, len(reqInfo.Headers))
		for key := range reqInfo.Headers {
			headerNames = append(headerNames, key)
		}
		sort.Strings(headerNames)
		headers := make([]string, len(reqInfo.Headers))
		for i := range headerNames {
			headers[i] = fmt.Sprintf("%s: %s", headerNames[i], strings.Join(reqInfo.Headers[headerNames[i]], ", "))
		}

		// append http request info fields
		l.logger = l.logger.With(
			"HTTP_RequestID", reqInfo.RequestID,
			reqInfo.Method, reqInfo.URL,
		)

		if l.showHeaders && len(headers) > 0 {
			l.logger = l.logger.With("HTTP_Headers", headers)
		}

		if l.showBody && len(reqInfo.Body) > 0 {
			l.logger = l.logger.With("HTTP_Body", string(reqInfo.Body))
		}
	}

	return l
}
