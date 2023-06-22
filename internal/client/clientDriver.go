package client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxiorder/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientDriver struct {
	driverClient proto.DriverServiceClient
	driverConn   *grpc.ClientConn

	cfg *config.Config
}

func NewClientDriver(cfg *config.Config) (*ClientDriver, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	driverConn, err := grpc.Dial(cfg.GRPC_DIVER_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}
	driverClient := proto.NewDriverServiceClient(driverConn)

	return &ClientDriver{
		driverClient: driverClient,
		driverConn:   driverConn,

		cfg: cfg,
	}, nil
}

func (c *ClientDriver) SyncDriver(ctx context.Context, drivers []*proto.Driver) ([]*proto.Driver, error) {
	request := &proto.Info{
		Drivers: drivers,
	}
	response, err := c.driverClient.SyncDriver(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("find driver failed: %w", err)
	}

	var syncDrivers []*proto.Driver
	for _, driver := range response.Drivers {
		d := driver
		syncDrivers = append(syncDrivers, d)
	}
	return syncDrivers, nil
}

func (c *ClientDriver) Close() error {
	return c.driverConn.Close()
}
