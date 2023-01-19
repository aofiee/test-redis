package services

import (
	"errors"
	"reflect"
	"taobin-service/internal/core/domain"
	"taobin-service/mocks"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/go-redis/redis/v9"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
)

type Suite struct {
	redis  *redis.Client
	reMock redismock.ClientMock
}

func NewSuite() *Suite {
	s := &Suite{}
	s.redis, s.reMock = redismock.NewClientMock()
	return s
}

func TestNewService(t *testing.T) {
	s := NewSuite()
	repo := new(mocks.Repository)
	service := New(repo, s.redis)
	assert.Equal(t, service, &Service{repo: repo, redis: s.redis})
}

func TestUpdateStock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		s := NewSuite()
		repo := new(mocks.Repository)
		service := New(repo, s.redis)

		request := domain.StockRequest{}
		err := faker.FakeData(&request)
		assert.NoError(t, err)
		returnResponse := domain.StockResponse{}
		repo.On("UpdateStock", request).Return(&returnResponse, nil)
		response, err := service.UpdateStock(request)
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, reflect.TypeOf(response), reflect.TypeOf(&domain.StockResponse{}))
	})
	t.Run("error", func(t *testing.T) {
		s := NewSuite()
		repo := new(mocks.Repository)
		service := New(repo, s.redis)

		request := domain.StockRequest{}
		err := faker.FakeData(&request)
		assert.NoError(t, err)
		returnResponse := domain.StockResponse{}
		repo.On("UpdateStock", request).Return(&returnResponse, errors.New("error"))
		response, err := service.UpdateStock(request)
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, reflect.TypeOf(response), reflect.TypeOf(&domain.StockResponse{}))
	})
}
