package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/RipperAcskt/innotaxiorder/internal/service"
	"github.com/RipperAcskt/innotaxiorder/internal/service/mocks"
	"github.com/bmizerany/assert"
	"github.com/golang/mock/gomock"
)

func TestCreateOrder(t *testing.T) {
	type mockBehavior func(s *mocks.MockRepo)
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "create order",
			mockBehavior: func(s *mocks.MockRepo) {
				s.EXPECT().CreateOrder(context.Background(), model.Order{}).Return(nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepo(ctrl)
			DriverRepo := mocks.NewMockDriverService(ctrl)
			DriverService := service.New(repo, DriverRepo, &config.Config{SYNC_TIME: 10})

			tt.mockBehavior(repo)

			service := service.Service{
				Repo:          repo,
				DriverService: DriverService,
			}

			err := service.CreateOrder(context.Background(), model.Order{})
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestGetOrders(t *testing.T) {
	type mockBehavior func(s *mocks.MockRepo)
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "get orders",
			mockBehavior: func(s *mocks.MockRepo) {
				s.EXPECT().GetOrders(context.Background(), nil).Return(nil, nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepo(ctrl)
			DriverRepo := mocks.NewMockDriverService(ctrl)
			DriverService := service.New(repo, DriverRepo, &config.Config{SYNC_TIME: 10})

			tt.mockBehavior(repo)

			service := service.Service{
				Repo:          repo,
				DriverService: DriverService,
			}

			_, err := service.GetOrder(context.Background(), nil)
			assert.Equal(t, err, tt.err)
		})
	}
}

func TestCancelOrder(t *testing.T) {
	type mockBehavior func(s *mocks.MockRepo, s1 *mocks.MockDriverService)
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "cancel order",
			mockBehavior: func(s *mocks.MockRepo, s1 *mocks.MockDriverService) {
				s.EXPECT().GetOrders(context.Background(), nil).Return(nil, nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepo(ctrl)
			DriverRepo := mocks.NewMockDriverService(ctrl)
			DriverService := service.New(repo, DriverRepo, &config.Config{SYNC_TIME: 10})

			tt.mockBehavior(repo, DriverRepo)

			service := service.Service{
				Repo:          repo,
				DriverService: DriverService,
			}

			_, err := service.CancelOrder(context.Background(), "")
			assert.NotEqual(t, err, tt.err)
		})
	}
}

func TestSetRating(t *testing.T) {
	type mockBehavior func(s *mocks.MockRepo, s1 *mocks.MockDriverService)
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "set raiting",
			mockBehavior: func(s *mocks.MockRepo, s1 *mocks.MockDriverService) {
				s.EXPECT().GetOrders(context.Background(), nil).Return(nil, nil)
			},
			err: nil,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepo(ctrl)
			DriverRepo := mocks.NewMockDriverService(ctrl)
			DriverService := service.New(repo, DriverRepo, &config.Config{SYNC_TIME: 10})

			tt.mockBehavior(repo, DriverRepo)

			service := service.Service{
				Repo:          repo,
				DriverService: DriverService,
			}

			_, err := service.SetRatingService(context.Background(), model.Rating{}, "0", "")
			assert.NotEqual(t, err, tt.err)
		})
	}
}

func TestFind(t *testing.T) {
	type mockBehavior func(s *mocks.MockRepo, s1 *mocks.MockDriverService)
	test := []struct {
		name         string
		mockBehavior mockBehavior
		err          error
	}{
		{
			name: "find driver",
			mockBehavior: func(s *mocks.MockRepo, s1 *mocks.MockDriverService) {
				s.EXPECT().GetOrdersByUserID(context.Background(), "", model.StatusWaiting.String()).Return(nil, nil)
			},
			err: fmt.Errorf("you haven't any waiting orders"),
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepo(ctrl)
			DriverRepo := mocks.NewMockDriverService(ctrl)
			DriverService := service.New(repo, DriverRepo, &config.Config{SYNC_TIME: 10})

			tt.mockBehavior(repo, DriverRepo)

			service := service.Service{
				Repo:          repo,
				DriverService: DriverService,
			}

			_, err := service.Find(context.Background(), "")
			assert.Equal(t, err, tt.err)
		})
	}
}
