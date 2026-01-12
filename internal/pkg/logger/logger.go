package logger

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	"boilerplate/internal/pkg/metadata"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	fieldTraceID   = "trace_id"
	fieldSpanID    = "span_id"
	fieldRequestID = "request_id"
	fieldUserID    = "user_id"
	fieldIP        = "ip"
)

type Logger interface {
	With(key, value string) Logger

	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, msg string, args ...any)
	DebugKV(ctx context.Context, msg string, values ...any)

	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, msg string, args ...any)
	InfoKV(ctx context.Context, msg string, values ...any)

	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, msg string, args ...any)
	WarnKV(ctx context.Context, msg string, values ...any)

	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, msg string, args ...any)
	ErrorKV(ctx context.Context, msg string, values ...any)

	Panic(ctx context.Context, msg string)
	Panicf(ctx context.Context, msg string, args ...any)
	PanicKV(ctx context.Context, msg string, values ...any)

	Fatal(ctx context.Context, msg string)
	Fatalf(ctx context.Context, msg string, args ...any)
	FatalKV(ctx context.Context, msg string, values ...any)

	IsWithDebug() bool

	Close() error
}

type logger struct {
	logger *zap.Logger
}

func New(opts ...ConfigOption) (Logger, error) {
	l, err := NewConfig(opts...).Build()
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}

	return &logger{
		logger: l,
	}, nil
}

func (l *logger) With(key, value string) Logger {
	return &logger{
		logger: l.logger.With(zap.String(key, value)),
	}
}

func (l *logger) Close() error {
	err := l.logger.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		return fmt.Errorf("logger sync: %w", err)
	}
	return nil
}

func (l *logger) Debug(ctx context.Context, msg string) {
	l.logger.Debug(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Debugf(ctx context.Context, msg string, args ...any) {
	l.logger.Debug(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) DebugKV(ctx context.Context, msg string, values ...any) {
	l.logger.Debug(msg, l.getFields(ctx, values)...)
}

func (l *logger) Info(ctx context.Context, msg string) {
	l.logger.Info(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Infof(ctx context.Context, msg string, args ...any) {
	l.logger.Info(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) InfoKV(ctx context.Context, msg string, values ...any) {
	l.logger.Info(msg, l.getFields(ctx, values)...)
}

func (l *logger) Warn(ctx context.Context, msg string) {
	l.logger.Warn(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Warnf(ctx context.Context, msg string, args ...any) {
	l.logger.Warn(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) WarnKV(ctx context.Context, msg string, values ...any) {
	l.logger.Warn(msg, l.getFields(ctx, values)...)
}

func (l *logger) Error(ctx context.Context, msg string) {
	l.logger.Error(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Errorf(ctx context.Context, msg string, args ...any) {
	l.logger.Error(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) ErrorKV(ctx context.Context, msg string, values ...any) {
	l.logger.Error(msg, l.getFields(ctx, values)...)
}

func (l *logger) Panic(ctx context.Context, msg string) {
	l.logger.Panic(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Panicf(ctx context.Context, msg string, args ...any) {
	l.logger.Panic(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) PanicKV(ctx context.Context, msg string, values ...any) {
	l.logger.Panic(msg, l.getFields(ctx, values)...)
}

func (l *logger) Fatal(ctx context.Context, msg string) {
	l.logger.Fatal(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Fatalf(ctx context.Context, msg string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) FatalKV(ctx context.Context, msg string, values ...any) {
	l.logger.Fatal(msg, l.getFields(ctx, values)...)
}

func (l *logger) getFields(ctx context.Context, values []any) []zap.Field {
	fields := make([]zap.Field, 0, len(values)/2+2)

	for i := 0; i < len(values)-1; i += 2 {
		key, ok := values[i].(string)
		if !ok {
			key = fmt.Sprint(values[i])
		}

		switch v := values[i+1].(type) {
		case string:
			fields = append(fields, zap.String(key, v))
		case int:
			fields = append(fields, zap.Int(key, v))
		case bool:
			fields = append(fields, zap.Bool(key, v))
		case float64:
			fields = append(fields, zap.Float64(key, v))
		default:
			fields = append(fields, zap.Any(key, v))
		}
	}

	traceID := l.getTraceID(ctx)
	if traceID != "" {
		fields = append(fields, zap.String(fieldTraceID, traceID))
	}

	spanID := l.getSpanID(ctx)
	if spanID != "" {
		fields = append(fields, zap.String(fieldSpanID, spanID))
	}

	requestID, exist := metadata.GetRequestID(ctx)
	if exist {
		fields = append(fields, zap.String(fieldRequestID, requestID))
	}

	userID, exist := metadata.GetUserID(ctx)
	if exist {
		fields = append(fields, zap.Int(fieldUserID, userID))
	}

	ip, exist := metadata.GetIP(ctx)
	if exist {
		fields = append(fields, zap.String(fieldIP, ip))
	}

	return fields
}

func (l *logger) getTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func (l *logger) getSpanID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}

func (l *logger) IsWithDebug() bool {
	return l.logger.Core().Enabled(zap.DebugLevel)
}
