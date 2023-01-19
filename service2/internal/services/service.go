package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"taobin-service/internal/core/domain"
	"taobin-service/internal/core/ports"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	StockWarningSubscriber      = "StockWarningSubscriber"
	MachineRefillSubscriber     = "MachineRefillSubscriber"
	StockLevelOkEvent           = "StockLevelOkSubscriber"
	StockRefillItems        int = 10
)

type Service struct {
	repo  ports.Repository
	redis *redis.Client
}

func New(repo ports.Repository, redis *redis.Client) *Service {
	service := &Service{
		repo:  repo,
		redis: redis,
	}
	service.SubscribeStockWarning(redis)
	return service
}

func (s *Service) SubscribeStockWarning(redis *redis.Client) {
	go func() {
		channel := []string{StockWarningSubscriber, StockLevelOkEvent}
		var machine domain.MachineResponse
		ctx := context.Background()
		pubsub := redis.Subscribe(ctx, channel...)
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				logrus.Errorln(err)
			}
			if msg.Channel == StockWarningSubscriber {
				err = json.Unmarshal([]byte(msg.Payload), &machine)
				if err != nil {
					logrus.Errorln(err)
				}
				fmt.Println(msg.Channel, msg.Payload)
				err = s.MachineRefillEvent(machine)
				if err != nil {
					logrus.Errorln(err)
				}
			}
			if msg.Channel == StockLevelOkEvent {
				fmt.Println(msg.Channel, msg.Payload)
			}
		}
	}()
}

func (s *Service) MachineRefillEvent(payload domain.MachineResponse) error {
	var (
		stockRemain int
		updatedBy   = "system"
	)
	query := domain.QueryStockRequest{}
	stockList, err := s.GetStock(query)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	if stockList.Stocks == nil || len(stockList.Stocks) == 0 {
		return errors.New("stock is empty")
	}
	stockID := int(*stockList.Stocks[0].ID)
	if stockList.Stocks[0].Stock != nil {
		stockRemain = *stockList.Stocks[0].Stock - StockRefillItems
		if stockRemain < 0 {
			stockRemain = 100000
		}
	}
	update := domain.StockRequest{
		ID:        &stockID,
		Stock:     &stockRemain,
		UpdatedBy: &updatedBy,
	}
	_, err = s.UpdateStock(update)
	if err != nil {
		logrus.Errorln(err)
		return err
	}
	*payload.Stock += StockRefillItems
	pubData, err := json.Marshal(payload)
	if err != nil {
		logrus.Errorln(err)
	}
	return s.redis.Publish(context.Background(), MachineRefillSubscriber, pubData).Err()
}

func (s *Service) UpdateStock(payload domain.StockRequest) (*domain.StockResponse, error) {
	stock, err := s.repo.UpdateStock(payload)
	if err != nil {
		return &domain.StockResponse{}, err
	}
	return stock, nil
}

func (s *Service) GetStock(payload domain.QueryStockRequest) (*domain.StockListResult, error) {
	var (
		page    int
		perPage int
		offset  int
	)

	if payload.Page != nil {
		page = *payload.Page
	} else {
		page = 1
		payload.Page = &page
	}
	if payload.Limit != nil {
		perPage = *payload.Limit
	} else {
		perPage = 100
		payload.Limit = &perPage
	}
	offset = (page - 1) * perPage

	payload.Pagination = &domain.Pagination{
		Limit:  perPage,
		Offset: offset,
	}

	if payload.OrderBy != nil {
		asc := true
		if payload.Asc != nil {
			asc = *payload.Asc
		}
		payload.SortMethod = &domain.SortMethod{
			Asc:     asc,
			OrderBy: *payload.OrderBy,
		}
	} else {
		payload.SortMethod = &domain.SortMethod{
			Asc:     true,
			OrderBy: "ID",
		}
	}
	return s.repo.GetStock(payload)
}
