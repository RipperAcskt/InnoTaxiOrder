package client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxiorder/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	analystClient proto.AnalystServiceClient
	analystConn   *grpc.ClientConn

	driverClient proto.DriverServiceClient
	driverConn   *grpc.ClientConn

	cfg *config.Config
}

type OrderRequest struct {
	Id       string
	TaxiType string
}

func New(cfg *config.Config) (*Clients, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	analystConn, err := grpc.Dial(cfg.GRPC_ANALYST_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}
	analystClient := proto.NewAnalystServiceClient(analystConn)

	driverConn, err := grpc.Dial(cfg.GRPC_DIVER_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}
	driverClient := proto.NewDriverServiceClient(driverConn)

	return &Clients{
		analystClient: analystClient,
		driverClient:  driverClient,

		analystConn: analystConn,
		driverConn:  driverConn,

		cfg: cfg,
	}, nil
}

func (c *Clients) SyncDriver(ctx context.Context, drivers []*proto.Driver) ([]*proto.Driver, error) {
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

func (c *Clients) SetRating(ctx context.Context, raiting *proto.Rating) error {
	_, err := c.analystClient.SetRating(ctx, raiting)
	if err != nil {
		return fmt.Errorf("set rating failed: %w", err)
	}
	return nil
}

func (c *Clients) Close() error {
	err := c.analystConn.Close()
	if err != nil {
		return fmt.Errorf("user conn close failed: %w", err)
	}

	err = c.driverConn.Close()
	if err != nil {
		return fmt.Errorf("driver conn close failed: %w", err)
	}
	return nil
}
