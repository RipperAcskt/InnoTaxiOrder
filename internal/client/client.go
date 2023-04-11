package client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
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

func (u *User) SyncDriver(ctx context.Context, order orderProto.Info) (*model.Order, error) {
	request := &orderProto.Info{
		OrderID:  order.Id,
		TaxiType: order.TaxiType,
	}
	response, err := u.client.SyncDriver(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("find driver failed: %w", err)
	}
	return &model.Order{DriverID: response.ID, DriverName: response.Name, DriverPhone: response.PhoneNumber, DriverRaiting: response.Raiting}, nil
}

func (u *User) Close() error {
	return u.conn.Close()
}
