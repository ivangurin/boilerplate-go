package logger

import (
	"context"
	"errors"
	"fmt"
	"syscall"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	fieldTraceID = "trace_id"
	fieldSpanID  = "span_id"
)

type Logger struct {
	logger *zap.Logger
}

func New(opts ...ConfigOption) (*Logger, error) {
	logger, err := NewConfig(opts...).Build()
	if err != nil {
		return nil, fmt.Errorf("create logger: %s", err.Error())
	}

	return &Logger{
		logger: logger,
	}, nil
}

func (l *Logger) Close() error {
	err := l.logger.Sync()
	if err != nil && !errors.Is(err, syscall.ENOTTY) {
		return fmt.Errorf("logger sync: %s", err.Error())
	}
	return nil
}

func (l *Logger) Debug(ctx context.Context, msg string) {
	l.logger.Debug(msg, l.getFields(ctx, nil)...)
}

func (l *Logger) Debugf(ctx context.Context, msg string, args ...any) {
	l.logger.Debug(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *Logger) DebugKV(ctx context.Context, msg string, values ...string) {
	l.logger.Debug(msg, l.getFields(ctx, values)...)
}

func (l *Logger) Info(ctx context.Context, msg string) {
	l.logger.Info(msg, l.getFields(ctx, nil)...)
}

func (l *Logger) Infof(ctx context.Context, msg string, args ...any) {
	l.logger.Info(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *Logger) InfoKV(ctx context.Context, msg string, values ...string) {
	l.logger.Info(msg, l.getFields(ctx, values)...)
}

func (l *Logger) Warn(ctx context.Context, msg string) {
	l.logger.Warn(msg, l.getFields(ctx, nil)...)
}

func (l *Logger) Warnf(ctx context.Context, msg string, args ...any) {
	l.logger.Warn(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *Logger) WarnKV(ctx context.Context, msg string, values ...string) {
	l.logger.Warn(msg, l.getFields(ctx, values)...)
}

func (l *Logger) Error(ctx context.Context, msg string) {
	l.logger.Error(msg, l.getFields(ctx, nil)...)
}

func (l *Logger) Errorf(ctx context.Context, msg string, args ...any) {
	l.logger.Error(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *Logger) ErrorKV(ctx context.Context, msg string, values ...string) {
	l.logger.Error(msg, l.getFields(ctx, values)...)
}

func (l *Logger) Panic(ctx context.Context, msg string) {
	l.logger.Panic(msg, l.getFields(ctx, nil)...)
}

func (l *Logger) Panicf(ctx context.Context, msg string, args ...any) {
	l.logger.Panic(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *Logger) PanicKV(ctx context.Context, msg string, values ...string) {
	l.logger.Panic(msg, l.getFields(ctx, values)...)
}

func (l *Logger) Fatal(ctx context.Context, msg string) {
	l.logger.Fatal(msg, l.getFields(ctx, nil)...)
}

func (l *Logger) Fatalf(ctx context.Context, msg string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *Logger) FatalKV(ctx context.Context, msg string, values ...string) {
	l.logger.Fatal(msg, l.getFields(ctx, values)...)
}

func (l *Logger) getFields(ctx context.Context, values []string) []zap.Field {
	fields := make([]zap.Field, 0, len(values)/2+2)

	for i := 0; i < len(values)-1; i += 2 {
		fields = append(fields, zap.String(values[i], fmt.Sprint(values[i+1])))
	}

	traceID := l.getTraceID(ctx)
	if traceID != "" {
		fields = append(fields, zap.String(fieldTraceID, traceID))
	}

	spanID := l.getSpanID(ctx)
	if spanID != "" {
		fields = append(fields, zap.String(fieldSpanID, spanID))
	}

	return fields
}

func (l *Logger) getTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func (l *Logger) getSpanID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}
