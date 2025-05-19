package jaeger_app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Endpoint    string
	ServiceName string
	Username    string
	Password    string
	CaCertPath  string
}

type JaegerApp struct {
	cfg            *Config
	tracerProvider *trace.TracerProvider
}

func New(cfg *Config) *JaegerApp {
	return &JaegerApp{cfg: cfg}
}

func (j *JaegerApp) Start(ctx context.Context) error {
	var opts []otlptracegrpc.Option

	opts = append(opts, otlptracegrpc.WithEndpoint(j.cfg.Endpoint))
	headers := map[string]string{}

	if j.cfg.Username != "" && j.cfg.Password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", j.cfg.Username, j.cfg.Password)))
		headers["Authorization"] = "Basic " + auth
	}

	opts = append(opts, otlptracegrpc.WithHeaders(headers))

	if j.cfg.CaCertPath != "" {
		caCert, err := os.ReadFile(j.cfg.CaCertPath)
		if err != nil {
			return fmt.Errorf("failed to read CA certificate: %w", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		creds := credentials.NewTLS(&tls.Config{RootCAs: caCertPool})
		opts = append(opts, otlptracegrpc.WithTLSCredentials(creds))
	} else {
		opts = append(opts, otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
	}

	exp, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return err
	}

	j.tracerProvider = trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(j.cfg.ServiceName),
		)),
	)

	otel.SetTracerProvider(j.tracerProvider)
	return nil
}

func (j *JaegerApp) Stop(ctx context.Context) error {
	if j.tracerProvider != nil {
		return j.tracerProvider.Shutdown(ctx)
	}

	return nil
}
