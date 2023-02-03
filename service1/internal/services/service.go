package services

import (
	"context"
	"encoding/json"
	"fmt"
	"taobin-service/internal/core/domain"
	"taobin-service/internal/core/ports"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	minimumStock            int = 3
	LowStockWarningEvent        = "StockWarningSubscriber"
	MachineRefillSubscriber     = "MachineRefillSubscriber"
	StockLevelOkEvent           = "StockLevelOkSubscriber"
	jobKey                      = "convert"
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
	service.Worker(redis)
	service.MachineRefillSubscriber(redis)
	return service
}

func (s *Service) Worker(redis *redis.Client) {
	go func() {
		logrus.Info("Starting worker")
		for {
			result, err := redis.BLPop(context.Background(), 0*time.Second, jobKey).Result()
			if err != nil {
				logrus.Error(err)
			}
			if len(result) == 2 {
				fmt.Println("Executing job: ", result[1])
			}
		}
	}()
}

func (s *Service) MachineRefillSubscriber(redis *redis.Client) {
	go func() {
		var machine domain.MachineResponse
		ctx := context.Background()
		pubsub := redis.Subscribe(ctx, MachineRefillSubscriber)
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				logrus.Errorln(err)
			}
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
	}()
}

func (s *Service) MachineRefillEvent(payload domain.MachineResponse) error {
	machineID := int(*payload.ID)
	update := domain.MachineRequest{
		ID:        &machineID,
		Machine:   payload.Machine,
		Stock:     payload.Stock,
		UpdatedBy: payload.UpdatedBy,
	}
	_, err := s.repo.UpdateMachine(update)
	if err != nil {
		return err
	}
	// StockLevelOkEvent
	ctx := context.Background()
	err = s.redis.Publish(ctx, StockLevelOkEvent, "OK").Err()
	if err != nil {
		logrus.Errorln(err)
	}
	return nil
}

func (s *Service) CreateMachine(payload domain.MachineRequest) (*domain.MachineResponse, error) {
	payload.UpdatedBy = payload.CreatedBy
	machine, err := s.repo.CreateMachine(payload)
	if err != nil {
		return &domain.MachineResponse{}, err
	}
	return machine, nil
}

func (s *Service) UpdateMachine(payload domain.MachineRequest) (*domain.MachineResponse, error) {
	if *payload.Stock < minimumStock {
		ctx := context.Background()
		logrus.Info("LowStockWarningEvent")
		pubData, err := json.Marshal(payload)
		if err != nil {
			logrus.Errorln(err)
		}
		err = s.redis.Publish(ctx, LowStockWarningEvent, pubData).Err()
		if err != nil {
			logrus.Errorln(err)
		}
		result, err := s.redis.RPush(context.Background(), jobKey, pubData).Result()
		if err != nil {
			logrus.Errorln(err)
		}
		logrus.Info("Example Job queued: ", result)
	}
	machine, err := s.repo.UpdateMachine(payload)
	if err != nil {
		return &domain.MachineResponse{}, err
	}
	return machine, nil
}

func (s *Service) DeleteMachine(payload domain.MachineRequest) (*domain.MachineResponse, error) {
	machine, err := s.repo.DeleteMachine(payload)
	if err != nil {
		return &domain.MachineResponse{}, err
	}
	return machine, nil
}

func (s *Service) GetMachine(payload domain.QueryMachineRequest) (*domain.MachineListResult, error) {
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
	return s.repo.GetMachine(payload)
}
