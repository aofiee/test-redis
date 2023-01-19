package repositories

import (
	"database/sql"
	"regexp"
	"taobin-service/internal/core/domain"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bxcodec/faker"
	"github.com/go-redis/redis/v9"
	"github.com/go-redis/redismock/v9"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	fakeLayoutDate = "2006-01-02 03:04:05 +0700 UTC"
)

type Suite struct {
	dbGorm *gorm.DB
	redis  *redis.Client
	dbMock sqlmock.Sqlmock
	reMock redismock.ClientMock
}

func NewSuite() *Suite {
	s := &Suite{}
	var (
		db  *sql.DB
		err error
	)
	s.redis, s.reMock = redismock.NewClientMock()

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

func TestUpdateCreateStock(t *testing.T) {
	t.Run("no parameter to updates", func(t *testing.T) {
		var (
			request domain.StockRequest
			err     error
		)
		s := NewSuite()
		postgres := NewPostgres(s.dbGorm, s.redis)

		_, err = postgres.UpdateStock(request)
		assert.Error(t, err)
	})
	t.Run("success", func(t *testing.T) {
		var (
			request domain.StockRequest
			err     error
		)
		s := NewSuite()
		postgres := NewPostgres(s.dbGorm, s.redis)
		err = faker.FakeData(&request)
		assert.NoError(t, err)
		s.dbMock.ExpectBegin()
		s.dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "stock" WHERE "id" = $1 AND "stock"."deleted_at" IS NULL ORDER BY "stock"."id" LIMIT 1`)).
			WithArgs(request.ID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(request.ID))
		s.dbMock.ExpectExec(regexp.QuoteMeta(`UPDATE "stock" SET "created_by"=$1,"stock"=$2,"updated_by"=$4 WHERE "id" = $5`)).
			WithArgs(request.CreatedBy, request.Stock, request.UpdatedBy, request.ID).
			WillReturnResult(sqlmock.NewResult(1, int64(*request.ID)))
		_, err = postgres.UpdateStock(request)
		assert.NoError(t, err)
	})
}
