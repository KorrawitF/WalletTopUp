package lib

import (
	"WalletTopUp/internal/config"
	"WalletTopUp/pkg/utils"
	"context"
	"log"
	"time"
)

const (
	Reset  = "\033[0m"  // Normal
	Red    = "\033[31m" // Error
	Green  = "\033[32m" // Info
	Yellow = "\033[33m" // Warn
	Gray   = "\033[37m" // Debug
)

type Logger interface {
	Info(ctx context.Context, meta Meta)
	Error(ctx context.Context, meta Meta)
	Warn(ctx context.Context, meta Meta)
	Debug(ctx context.Context, meta Meta)
}

type logger struct {
	debugMode bool
	log       *log.Logger
}

func NewLogger(conf *config.Config, log *log.Logger) Logger {
	log.SetFlags(0)
	return &logger{
		debugMode: conf.Server.Mode != "release",
		log:       log,
	}
}

type Meta struct {
	Event string
	Msg   string
	Error error
}

func (l *logger) Info(ctx context.Context, meta Meta) {
	ts := time.Now().Format(time.RFC3339Nano)
	l.log.Printf(Green+"%s - [Info] - Event: %s - Tid: %s - Message: %s\n"+Reset, ts, meta.Event, utils.GetTraceId(ctx), meta.Msg)
}

func (l *logger) Error(ctx context.Context, meta Meta) {
	ts := time.Now().Format(time.RFC3339Nano)
	l.log.Printf(Red+"%s - [Error] - Event: %s - Tid: %s - Message: %s, Error: %v\n"+Reset, ts, meta.Event, utils.GetTraceId(ctx), meta.Msg, meta.Error)
}

func (l *logger) Warn(ctx context.Context, meta Meta) {
	ts := time.Now().Format(time.RFC3339Nano)
	l.log.Printf(Yellow+"%s - [Warn] - Event: %s - Tid: %s - Message: %s\n"+Reset, ts, meta.Event, utils.GetTraceId(ctx), meta.Msg)
}

func (l *logger) Debug(ctx context.Context, meta Meta) {
	if !l.debugMode {
		l.log.Print(l.debugMode)
		return
	}
	ts := time.Now().Format(time.RFC3339Nano)
	l.log.Printf(Gray+"%s - [DEBUG] - Event: %s - Tid: %s - Message: %s, Error: %v\n"+Reset, ts, meta.Event, utils.GetTraceId(ctx), meta.Msg, meta.Error)
}
