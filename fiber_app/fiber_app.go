package fiber_app

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type FiberApp struct {
	Host        string
	App         *fiber.App
	middlewares []fiber.Handler
}

func New(
	host string,
	middlewares ...fiber.Handler,
) *FiberApp {
	return &FiberApp{
		Host:        host,
		App:         fiber.New(fiber.Config{DisableStartupMessage: true}),
		middlewares: middlewares,
	}
}

func (a *FiberApp) Start(ctx context.Context) error {
	for _, middleware := range a.middlewares {
		a.App.Use(middleware)
	}

	a.App.Get("/health_check", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	go func() {
		if err := a.App.Listen(a.Host); err != nil {
			log.Fatalf("error occurred while running http server: %s", err)
		}
	}()

	return nil
}

func (a *FiberApp) Stop(ctx context.Context) error {
	return a.App.ShutdownWithContext(ctx)
}
