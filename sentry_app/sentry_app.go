package sentry_app

import (
	"context"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

type Config struct {
	Dsn              string
	TracesSampleRate float64
	SampleRate       float64
	IgnoreFrame      string
	FlushTimeout     time.Duration
}

type SentryApp struct {
	cfg *Config
}

func New(cfg *Config) *SentryApp {
	return &SentryApp{cfg: cfg}
}

func (s *SentryApp) Start(ctx context.Context) (err error) {
	return sentry.Init(
		sentry.ClientOptions{
			Dsn:              s.cfg.Dsn,
			TracesSampleRate: s.cfg.TracesSampleRate,
			SampleRate:       s.cfg.SampleRate,
			BeforeSend: func(event *sentry.Event, _ *sentry.EventHint) *sentry.Event {
				return filterAlertWrapper(event, s.cfg.IgnoreFrame)
			},
		},
	)
}

func (s *SentryApp) Stop(ctx context.Context) (err error) {
	sentry.Flush(s.cfg.FlushTimeout)

	return nil
}

func filterAlertWrapperFrames(frames []sentry.Frame, ignoreFrame string) []sentry.Frame {
	filteredFrames := make([]sentry.Frame, 0, len(frames))

	for i := range frames {
		frame := frames[i]

		if strings.Contains(frame.Module, ignoreFrame) {
			continue
		}

		filteredFrames = append(filteredFrames, frame)
	}

	return filteredFrames
}

func filterAlertWrapper(event *sentry.Event, ignoreFrame string) *sentry.Event {
	for _, ex := range event.Exception {
		if ex.Stacktrace == nil {
			continue
		}

		ex.Stacktrace.Frames = filterAlertWrapperFrames(ex.Stacktrace.Frames, ignoreFrame)
	}

	for _, th := range event.Threads {
		if th.Stacktrace == nil {
			continue
		}

		th.Stacktrace.Frames = filterAlertWrapperFrames(th.Stacktrace.Frames, ignoreFrame)
	}

	return event
}
