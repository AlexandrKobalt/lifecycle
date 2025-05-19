package grpc_server_app

import (
	"context"
	"net"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Config struct {
	Host              string        `validate:"required"`
	MaxConnectionIdle time.Duration `validate:"required"`
	Timeout           time.Duration `validate:"required"`
	MaxConnectionAge  time.Duration `validate:"required"`
	Time              time.Duration `validate:"required"`
}

type GrpcServerApp struct {
	cfg    *Config
	Server *grpc.Server
}

func New(cfg *Config) (*GrpcServerApp, error) {
	server := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionIdle: cfg.MaxConnectionIdle,
				Timeout:           cfg.Timeout,
				MaxConnectionAge:  cfg.MaxConnectionAge,
				Time:              cfg.Time,
			},
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	return &GrpcServerApp{cfg: cfg, Server: server}, nil
}

func (g *GrpcServerApp) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", g.cfg.Host)
	if err != nil {
		return err
	}

	go func() {
		if err := g.Server.Serve(listener); err != nil {
			log.Fatal().Err(err).Msg("failed to start grpc server")
		}
	}()

	return nil
}

func (g *GrpcServerApp) Stop(ctx context.Context) error {
	g.Server.GracefulStop()

	return nil
}
