package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/client"
	"github.com/RipperAcskt/innotaxiorder/internal/handler/graph"
	"github.com/RipperAcskt/innotaxiorder/internal/handler/grpc"
	"github.com/RipperAcskt/innotaxiorder/internal/repo/elastic"
	"github.com/RipperAcskt/innotaxiorder/internal/server"
	"github.com/RipperAcskt/innotaxiorder/internal/service"

	"go.uber.org/zap"
)

func Run() error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("config new failed: %w", err)
	}

	log, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("new production failed: %w", err)
	}

	repo, err := elastic.New(cfg)
	if err != nil {
		return fmt.Errorf("elastic new failed: %w", err)
	}

	client, err := client.New(cfg)
	if err != nil {
		return fmt.Errorf("client new failed: %v", err)
	}

	service := service.New(repo, client, cfg)

	go func() {
		err := <-service.Err
		log.Error("error: ", zap.Error(err))
	}()

	err = service.SyncDrivers(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("sync drivers failed: %w", err)
	}

	handler, err := graph.New(service, cfg, log)
	if err != nil {
		return fmt.Errorf("handler new failed: %w", err)
	}

	server := &server.Server{
		Log: log,
	}

	go func() {
		if err := server.Run(handler.InitRouters(), cfg); err != nil && err != http.ErrServerClosed {
			log.Error(fmt.Sprintf("server run failed: %v", err))
			return
		}
	}()

	grpcServer := grpc.New(log, service, cfg)
	go func() {
		if err := grpcServer.Run(); err != nil {
			log.Error(fmt.Sprintf("grpc server run failed: %v", err))
			return
		}
	}()

	if err := server.ShutDown(); err != nil {
		return fmt.Errorf("server shut down failed: %w", err)
	}
	return nil
}
