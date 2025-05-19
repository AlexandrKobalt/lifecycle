package telebotV4_app

import (
	"context"

	"gopkg.in/telebot.v4"
)

type TelebotApp struct {
	Bot         *telebot.Bot
	middlewares []telebot.MiddlewareFunc
}

func New(
	settings telebot.Settings,
	middlewares ...telebot.MiddlewareFunc,
) (*TelebotApp, error) {
	bot, err := telebot.NewBot(settings)
	if err != nil {
		return nil, err
	}

	return &TelebotApp{
		Bot:         bot,
		middlewares: middlewares,
	}, nil
}

func (a *TelebotApp) Start(ctx context.Context) error {
	for _, middleware := range a.middlewares {
		a.Bot.Use(middleware)
	}

	go a.Bot.Start()

	return nil
}

func (a *TelebotApp) Stop(ctx context.Context) error {
	a.Bot.Stop()

	return nil
}
