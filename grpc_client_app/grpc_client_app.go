package grpc_client_app

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Host     string `validate:"required"`
	CertFile string
}

type GrpcClientApp struct {
	cfg        *Config
	Connection *grpc.ClientConn
}

func New(cfg *Config) *GrpcClientApp {
	return &GrpcClientApp{cfg: cfg}
}

func (g *GrpcClientApp) Start(ctx context.Context) error {
	if g.Connection != nil {
		return nil
	}

	var opts []grpc.DialOption

	if g.cfg.CertFile != "" {
		creds, err := credentials.NewClientTLSFromFile(g.cfg.CertFile, "")
		if err != nil {
			return err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	connection, err := grpc.DialContext(ctx, g.cfg.Host, opts...)
	if err != nil {
		return err
	}

	g.Connection = connection
	return nil
}

func (g *GrpcClientApp) Stop(ctx context.Context) error {
	return g.Connection.Close()
}
