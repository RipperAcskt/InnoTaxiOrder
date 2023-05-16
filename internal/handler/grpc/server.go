package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/RipperAcskt/innotaxiorder/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	listener   net.Listener
	grpcServer *grpc.Server
	log        *zap.Logger
	service    *service.Service
	cfg        *config.Config
}

func New(log *zap.Logger, s *service.Service, cfg *config.Config) *Server {
	return &Server{nil, nil, log, s, cfg}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.cfg.GRPC_HOST)

	if err != nil {
		return fmt.Errorf("listen failed: %w", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	s.listener = listener
	s.grpcServer = grpcServer

	proto.RegisterOrderServiceServer(grpcServer, s)
	err = grpcServer.Serve(listener)
	if err != nil {
		return fmt.Errorf("serve failed: %w", err)
	}

	return nil
}

func (s *Server) GetOrderQuantity(c context.Context, params *proto.Time) (*proto.Amount, error) {
	filters := model.OrderFilters{
		FromDate: params.TimeStarted,
	}
	pagginationInfo := model.PagginationInfo{
		PagginationFlag: false,
	}
	orders, err := s.service.GetOrdersList(c, filters, pagginationInfo)
	if err != nil {
		return nil, fmt.Errorf("get orders list failed")
	}
	resp := &proto.Amount{
		NumberOfOrders: int64(len(orders)),
	}
	return resp, nil
}

func (s *Server) Stop() error {
	s.log.Info("Shuttig down grpc...")

	err := s.listener.Close()
	if err != nil {
		return fmt.Errorf("listener close failed: %w", err)
	}

	s.grpcServer.Stop()
	s.log.Info("Grpc server exiting.")
	return nil
}
