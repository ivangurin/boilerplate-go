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

type Logger interface {
	Debug(ctx context.Context, msg string)
	Debugf(ctx context.Context, msg string, args ...any)
	DebugKV(ctx context.Context, msg string, values ...string)

	Info(ctx context.Context, msg string)
	Infof(ctx context.Context, msg string, args ...any)
	InfoKV(ctx context.Context, msg string, values ...string)

	Warn(ctx context.Context, msg string)
	Warnf(ctx context.Context, msg string, args ...any)
	WarnKV(ctx context.Context, msg string, values ...string)

	Error(ctx context.Context, msg string)
	Errorf(ctx context.Context, msg string, args ...any)
	ErrorKV(ctx context.Context, msg string, values ...string)

	Panic(ctx context.Context, msg string)
	Panicf(ctx context.Context, msg string, args ...any)
	PanicKV(ctx context.Context, msg string, values ...string)

	Fatal(ctx context.Context, msg string)
	Fatalf(ctx context.Context, msg string, args ...any)
	FatalKV(ctx context.Context, msg string, values ...string)

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

func (l *logger) DebugKV(ctx context.Context, msg string, values ...string) {
	l.logger.Debug(msg, l.getFields(ctx, values)...)
}

func (l *logger) Info(ctx context.Context, msg string) {
	l.logger.Info(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Infof(ctx context.Context, msg string, args ...any) {
	l.logger.Info(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) InfoKV(ctx context.Context, msg string, values ...string) {
	l.logger.Info(msg, l.getFields(ctx, values)...)
}

func (l *logger) Warn(ctx context.Context, msg string) {
	l.logger.Warn(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Warnf(ctx context.Context, msg string, args ...any) {
	l.logger.Warn(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) WarnKV(ctx context.Context, msg string, values ...string) {
	l.logger.Warn(msg, l.getFields(ctx, values)...)
}

func (l *logger) Error(ctx context.Context, msg string) {
	l.logger.Error(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Errorf(ctx context.Context, msg string, args ...any) {
	l.logger.Error(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) ErrorKV(ctx context.Context, msg string, values ...string) {
	l.logger.Error(msg, l.getFields(ctx, values)...)
}

func (l *logger) Panic(ctx context.Context, msg string) {
	l.logger.Panic(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Panicf(ctx context.Context, msg string, args ...any) {
	l.logger.Panic(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) PanicKV(ctx context.Context, msg string, values ...string) {
	l.logger.Panic(msg, l.getFields(ctx, values)...)
}

func (l *logger) Fatal(ctx context.Context, msg string) {
	l.logger.Fatal(msg, l.getFields(ctx, nil)...)
}

func (l *logger) Fatalf(ctx context.Context, msg string, args ...any) {
	l.logger.Fatal(fmt.Sprintf(msg, args...), l.getFields(ctx, nil)...)
}

func (l *logger) FatalKV(ctx context.Context, msg string, values ...string) {
	l.logger.Fatal(msg, l.getFields(ctx, values)...)
}

func (l *logger) getFields(ctx context.Context, values []string) []zap.Field {
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
