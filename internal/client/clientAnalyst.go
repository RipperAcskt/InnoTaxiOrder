package client

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxiorder/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientAnalyst struct {
	analystClient proto.AnalystServiceClient
	analystConn   *grpc.ClientConn

	cfg *config.Config
}

type OrderRequest struct {
	Id       string
	TaxiType string
}

func NewClientAnalyst(cfg *config.Config) (*ClientAnalyst, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	analystConn, err := grpc.Dial(cfg.GRPC_ANALYST_SERVICE_HOST, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}
	analystClient := proto.NewAnalystServiceClient(analystConn)

	return &ClientAnalyst{
		analystClient: analystClient,
		analystConn:   analystConn,

		cfg: cfg,
	}, nil
}

func (c *ClientAnalyst) SetRating(ctx context.Context, raiting *proto.Rating) error {
	_, err := c.analystClient.SetRating(ctx, raiting)
	if err != nil {
		return fmt.Errorf("set rating failed: %w", err)
	}
	return nil
}

func (c *ClientAnalyst) Close() error {
	return c.analystConn.Close()
}
