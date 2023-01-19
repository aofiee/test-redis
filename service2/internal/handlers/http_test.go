package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"reflect"
	"taobin-service/internal/core/domain"
	"taobin-service/mocks"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type config struct {
	Env string
}

type Suite struct {
	dbGorm *gorm.DB
	dbMock sqlmock.Sqlmock
}

func NewSuite() *Suite {
	s := &Suite{}
	var (
		db  *sql.DB
		err error
	)
	db, s.dbMock, err = sqlmock.New()
	if err != nil {
		logrus.Error(err)
	}
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})
	s.dbGorm, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		logrus.Error(err)
	}
	return s
}

func TestNew(t *testing.T) {
	t.Run("should return a new instance of the handler", func(t *testing.T) {
		s := NewSuite()
		service := new(mocks.Service)
		handler := New(service, s.dbGorm)
		assert.Equal(t, reflect.TypeOf(handler), reflect.TypeOf(&HTTPHandler{}))
	})
}

func TestHealthCheck(t *testing.T) {
	t.Run("should return a 200 status code", func(t *testing.T) {
		s := NewSuite()
		service := new(mocks.Service)
		handler := New(service, s.dbGorm)

		app := fiber.New()
		app.Get("/test", handler.TestCheck)
		req, err := http.NewRequest("GET", "/test", nil)
		assert.NoError(t, err)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "completed")
	})
}

func TestUpdateStock(t *testing.T) {
	endpoint := "/"
	t.Run("should return a 200 status code", func(t *testing.T) {
		s := NewSuite()
		service := new(mocks.Service)
		handler := New(service, s.dbGorm)
		request := domain.StockRequest{}
		err := faker.FakeData(&request)
		assert.NoError(t, err)

		service.On("UpdateStock", mock.AnythingOfType(reflect.TypeOf(domain.StockRequest{}).String())).Return(&domain.StockResponse{}, nil)

		data, err := json.Marshal(request)
		if err != nil {
			t.Error(err)
		}
		payload := bytes.NewReader(data)

		app := fiber.New()
		app.Put(endpoint, handler.UpdateStock)
		req, err := http.NewRequest("PUT", endpoint, payload)
		assert.NoError(t, err)
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode, "completed")
	})
}
