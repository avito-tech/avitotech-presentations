package main

import (
	"context"
	"log/slog"
	"os"
)

// -----------------------------------------------

type HandlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddlware {
	return &HandlerMiddlware{next: next}
}

func (h *HandlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		if c.UserID != 0 {
			rec.Add("userID", c.UserID)
		}
		if c.Phone != "" {
			rec.Add("phone", c.Phone)
		}
		if c.Gate != "" {
			rec.Add("sms_gate", c.Gate)
		}
		if c.Message != "" {
			rec.Add("message", c.Message)
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}

type logCtx struct {
	UserID  int
	Phone   string
	Gate    string
	Message string
}

type keyType int

const key = keyType(0)

func WithLogUserID(ctx context.Context, userID int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.UserID = userID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{UserID: userID})
}

func WithLogPhone(ctx context.Context, phone string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Phone = phone
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Phone: phone})
}

func WithLogGate(ctx context.Context, gate string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Gate = gate
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Gate: gate})
}

func WithLogMessage(ctx context.Context, message string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Message = message
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Message: message})
}

// -----------------------------------------------

func InitLogging() {
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}

// -----------------------------------------------

func TransmitSMS(ctx context.Context, gate, phone, message string) error {
	slog.InfoContext(ctx, "Transmit SMS gateway OK")
	return nil
}

// -----------------------------------------------

func ResolveGate(ctx context.Context, phone string) (string, error) {
	gate := "RHO"
	slog.InfoContext(ctx, "Resolve SMS gateway OK")
	return gate, nil
}

// -----------------------------------------------

func SendSMS(ctx context.Context, phone string) error {
	slog.InfoContext(ctx, "Send SMS OK")
	message := "Спасибо"
	ctx = WithLogMessage(ctx, message)
	gate, _ := ResolveGate(ctx, phone)
	ctx = WithLogGate(ctx, gate)
	_ = TransmitSMS(ctx, gate, phone, message)
	return nil
}

// -----------------------------------------------

func GetPhoenByID(ctx context.Context, userID int) (string, error) {
	phone := "+78880001100"
	ctx = WithLogPhone(ctx, phone)
	slog.InfoContext(ctx, "Get phone OK")
	return phone, nil
}

// -----------------------------------------------

func Handler(ctx context.Context, userID int) {
	ctx = WithLogUserID(ctx, userID)
	slog.InfoContext(ctx, "Handler started")
	phone, _ := GetPhoenByID(ctx, userID)
	ctx = WithLogPhone(ctx, phone)
	_ = SendSMS(ctx, phone)
	slog.InfoContext(ctx, "Handler done")
}

// -----------------------------------------------

func main() {
	InitLogging()

	ctx := context.Background()
	Handler(ctx, 111)
}
