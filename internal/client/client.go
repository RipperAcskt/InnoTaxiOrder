package client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxiorder/config"
	orderProto "github.com/RipperAcskt/innotaxiorder/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type User struct {
	client orderProto.OrderServiceClient
	conn   *grpc.ClientConn
	cfg    *config.Config
}

type OrderRequest struct {
	Id       string
	TaxiType string
}

func New(cfg *config.Config) (*User, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial(cfg.GRPC_USER_SERVICE_HOST, opts...)

	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	client := orderProto.NewOrderServiceClient(conn)

	return &User{client, conn, cfg}, nil
}

func (u *User) SyncDriver(ctx context.Context, drivers []*orderProto.Driver) ([]*orderProto.Driver, error) {
	request := &orderProto.Info{
		Drivers: drivers,
	}
	response, err := u.client.SyncDriver(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("find driver failed: %w", err)
	}

	var syncDrivers []*orderProto.Driver
	for _, driver := range response.Drivers {
		d := driver
		syncDrivers = append(syncDrivers, d)
	}
	return syncDrivers, nil
}

func (u *User) Close() error {
	return u.conn.Close()
}
