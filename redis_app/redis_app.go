package redis_app

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host               string `validate:"required"`
	Port               string `validate:"required"`
	MinIdleConns       int    `validate:"required"`
	PoolSize           int    `validate:"required"`
	PoolTimeout        int    `validate:"required"`
	Password           string
	UseCertificates    bool
	InsecureSkipVerify bool
	CertificatesPaths  struct {
		Cert string
		Key  string
		Ca   string
	}
	DB int
}

type RedisApp struct {
	cfg    *Config
	Client *redis.Client
}

func New(cfg *Config) *RedisApp {
	return &RedisApp{cfg: cfg}
}

func (r *RedisApp) Start(ctx context.Context) error {
	opts := &redis.Options{}
	if r.cfg.UseCertificates {
		certs := make([]tls.Certificate, 0)
		if r.cfg.CertificatesPaths.Cert != "" && r.cfg.CertificatesPaths.Key != "" {
			cert, err := tls.LoadX509KeyPair(r.cfg.CertificatesPaths.Cert, r.cfg.CertificatesPaths.Key)
			if err != nil {
				return errors.Wrapf(
					err,
					"certPath: %v, keyPath: %v",
					r.cfg.CertificatesPaths.Cert,
					r.cfg.CertificatesPaths.Key,
				)
			}
			certs = append(certs, cert)
		}
		caCert, err := os.ReadFile(r.cfg.CertificatesPaths.Ca)
		if err != nil {
			return errors.Wrapf(err, "ca load path: %v", r.cfg.CertificatesPaths.Ca)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		opts = &redis.Options{
			Addr:         fmt.Sprintf("%s:%s", r.cfg.Host, r.cfg.Port),
			MinIdleConns: r.cfg.MinIdleConns,
			PoolSize:     r.cfg.PoolSize,
			PoolTimeout:  time.Duration(r.cfg.PoolTimeout) * time.Second,
			Password:     r.cfg.Password,
			DB:           r.cfg.DB,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: r.cfg.InsecureSkipVerify,
				Certificates:       certs,
				RootCAs:            caCertPool,
			},
		}
	} else {
		opts = &redis.Options{
			Addr:         fmt.Sprintf("%s:%s", r.cfg.Host, r.cfg.Port),
			MinIdleConns: r.cfg.MinIdleConns,
			PoolSize:     r.cfg.PoolSize,
			PoolTimeout:  time.Duration(r.cfg.PoolTimeout) * time.Second,
			Password:     r.cfg.Password,
			DB:           r.cfg.DB,
		}
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return errors.Wrapf(err, "ping")
	}

	r.Client = client

	return nil
}

func (r *RedisApp) Stop(ctx context.Context) error {
	return r.Client.Close()
}
