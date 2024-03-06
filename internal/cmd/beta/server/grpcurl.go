package server

import (
	"context"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/stateful/runme/v3/internal/config"
	runmetls "github.com/stateful/runme/v3/internal/tls"
)

func getDescSource(ctx context.Context, cfg *config.Config) (grpcurl.DescriptorSource, error) {
	cc, err := dialServer(ctx, cfg)
	if err != nil {
		return nil, err
	}
	client := grpcreflect.NewClientAuto(ctx, cc)
	return grpcurl.DescriptorSourceFromServer(ctx, client), nil
}

func dialServer(ctx context.Context, cfg *config.Config) (*grpc.ClientConn, error) {
	tlsConf, err := runmetls.LoadServerConfig(cfg.ServerTLSCertFile, cfg.ServerTLSKeyFile)
	if err != nil {
		return nil, err
	}
	// TODO(adamb): remove insecure skip verify
	tlsConf.InsecureSkipVerify = true

	creds := credentials.NewTLS(tlsConf)

	addr := cfg.ServerAddress
	network := "tcp"
	if strings.HasPrefix(addr, "unix://") {
		network = "unix"
		addr = strings.TrimPrefix(addr, "unix://")
	}

	var opts []grpc.DialOption

	cc, err := grpcurl.BlockingDial(ctx, network, addr, creds, opts...)
	return cc, errors.WithStack(err)
}
