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
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}

// -----------------------------------------------

func InitLogging() {
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}

// -----------------------------------------------

func TransmitSMS(ctx context.Context, gate, phone, message string) error {
	slog.InfoContext(ctx, "Transmit SMS gateway OK", "phone", phone, "sms_gate", gate, "message", message)
	return nil
}

// -----------------------------------------------

func ResolveGate(ctx context.Context, phone string) (string, error) {
	gate := "RHO"
	slog.InfoContext(ctx, "Resolve SMS gateway OK", "phone", phone, "sms_gate", gate)
	return gate, nil
}

// -----------------------------------------------

func SendSMS(ctx context.Context, phone string) error {
	slog.InfoContext(ctx, "Sending SMS", "phone", phone)
	message := "Спасибо"
	gate, _ := ResolveGate(ctx, phone)
	_ = TransmitSMS(ctx, gate, phone, message)
	slog.InfoContext(ctx, "Send SMS OK", "phone", phone, "message", message)
	return nil
}

// -----------------------------------------------

func GetPhoenByID(ctx context.Context, userID int) (string, error) {
	phone := "+78880001100"
	slog.InfoContext(ctx, "Get phone OK", "userID", userID, "phone", phone)
	return phone, nil
}

// -----------------------------------------------

func Handler(ctx context.Context, userID int) {
	slog.InfoContext(ctx, "Handler started", "userID", userID)
	phone, _ := GetPhoenByID(ctx, userID)
	_ = SendSMS(ctx, phone)
	slog.InfoContext(ctx, "Handler done", "userID", userID)
}

// -----------------------------------------------

func main() {
	InitLogging()

	// где-то в listener-e
	ctx := context.Background()
	Handler(ctx, 111)
}
