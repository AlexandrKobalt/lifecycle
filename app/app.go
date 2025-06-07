package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AlexandrKobalt/lifecycle"
)

type App interface {
	lifecycle.Lifecycle
	GetStartTimeout() time.Duration
	GetStopTimeout() time.Duration
}

func Run(a App) (err error) {
	startCtx, startCancel := context.WithTimeout(context.Background(), a.GetStartTimeout())
	defer startCancel()

	if err := a.Start(startCtx); err != nil {
		return err
	}

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-quitCh

	stopCtx, stopCancel := context.WithTimeout(context.Background(), a.GetStopTimeout())
	defer stopCancel()

	return a.Stop(stopCtx)
}
