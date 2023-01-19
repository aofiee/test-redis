package repositories

import (
	"errors"
	"taobin-service/internal/core/domain"
	"time"

	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

const (
	layoutDate = "2006-01-02"
)

type Postgres struct {
	dbGorm *gorm.DB
	Redis  *redis.Client
}

func NewPostgres(dbGorm *gorm.DB, redis *redis.Client) *Postgres {
	domain.MigrateDatabase(dbGorm)
	return &Postgres{
		dbGorm: dbGorm,
		Redis:  redis,
	}
}

func (r *Postgres) UpdateStock(payload domain.StockRequest) (*domain.StockResponse, error) {
	var (
		response domain.StockResponse
		result   domain.Stock
	)
	query := domain.QueryStockRequest{
		ID: payload.ID,
	}
	condition := r.machineCondition(query)
	columns := r.updateStockColumns(payload)
	tx := r.dbGorm.Begin()
	tx.Table(domain.Stock{}.TableName()).Where(condition).Updates(columns)
	if tx.Error != nil {
		tx.Rollback()
		return &domain.StockResponse{}, tx.Error
	}
	tx.Where(condition).First(&result)
	if tx.Error != nil {
		tx.Rollback()
		return &domain.StockResponse{}, tx.Error
	}
	tx.Commit()
	if result.ID == 0 {
		return &domain.StockResponse{}, errors.New("Not found")
	}
	response = domain.StockResponse{
		ID:        &result.ID,
		Stock:     result.Stock,
		CreatedBy: result.CreatedBy,
		UpdatedBy: result.UpdatedBy,
		DeletedBy: result.DeletedBy,
		CreatedAt: &result.CreatedAt,
		UpdatedAt: &result.UpdatedAt,
		DeletedAt: &result.DeletedAt,
	}
	return &response, nil
}

func (r *Postgres) machineCondition(condition domain.QueryStockRequest) map[string]interface{} {
	expression := map[string]interface{}{}
	if condition.ID != nil {
		expression["id"] = condition.ID
	}
	if condition.Stock != nil {
		expression["stock"] = condition.Stock
	}
	if condition.CreatedBy != nil {
		expression["created_by"] = condition.CreatedBy
	}
	if condition.UpdatedBy != nil {
		expression["updated_by"] = condition.UpdatedBy
	}
	if condition.DeletedBy != nil {
		expression["deleted_by"] = condition.DeletedBy
	}
	return expression
}

func (r *Postgres) updateStockColumns(payload domain.StockRequest) map[string]interface{} {
	expression := map[string]interface{}{}
	if payload.Stock != nil {
		expression["stock"] = payload.Stock
	}
	if payload.UpdatedBy != nil {
		expression["updated_by"] = payload.UpdatedBy
	}
	if payload.DeletedBy != nil {
		expression["deleted_by"] = payload.DeletedBy
	}
	return expression
}

func (r *Postgres) GetStock(payload domain.QueryStockRequest) (*domain.StockListResult, error) {
	var (
		machines []domain.Stock
	)
	condition := r.machineCondition(payload)
	tx := r.dbGorm
	count := r.dbGorm
	if payload.DeletedBy != nil || payload.DeletedAt != nil {
		tx = tx.Unscoped().Where(condition)
		count = count.Unscoped().Select("COUNT(id)").Where(condition)
	} else {
		tx = tx.Where(condition)
		count = count.Model(&domain.Stock{}).Select("COUNT(id)").Where(condition)
	}
	if payload.ID == nil {
		var order string
		if payload.SortMethod.OrderBy != "" {
			order = payload.SortMethod.OrderBy
		} else {
			order = "id"
		}
		if payload.SortMethod.Asc {
			tx = tx.Order(order + " ASC")
		} else {
			tx = tx.Order(order + " DESC")
		}
		tx = tx.Limit(payload.Pagination.Limit).Offset(payload.Pagination.Offset)
	}
	if payload.CreatedAt != nil {
		_, error := time.Parse(layoutDate, *payload.CreatedAt)
		if error != nil {
			return nil, error
		}
		tx = tx.Where("DATE(created_at) = ?", *payload.CreatedAt)
		count = count.Where("DATE(created_at) = ?", *payload.CreatedAt)
	}
	if payload.UpdatedAt != nil {
		_, error := time.Parse(layoutDate, *payload.UpdatedAt)
		if error != nil {
			return nil, error
		}
		tx = tx.Where("DATE(updated_at) = ?", *payload.UpdatedAt)
		count = count.Where("DATE(updated_at) = ?", *payload.UpdatedAt)
	}
	if payload.DeletedAt != nil {
		_, error := time.Parse(layoutDate, *payload.DeletedAt)
		if error != nil {
			return nil, error
		}
		tx = tx.Where("DATE(deleted_at) = ?", *payload.DeletedAt)
		count = count.Where("DATE(deleted_at) = ?", *payload.DeletedAt)
	}
	if len(payload.IDs) > 0 {
		tx.Where("id IN (?)", payload.IDs)
		count.Where("id IN (?)", payload.IDs)
	}
	tx.Find(&machines)
	if tx.Error != nil {
		return nil, tx.Error
	}
	result := domain.StockListResult{
		Stocks: []domain.StockResponse{},
	}
	result.CurrentPage = payload.Page
	result.PerPage = &payload.Pagination.Limit
	var toatalItem int64
	count.Count(&toatalItem)
	result.TotalItem = &toatalItem
	for _, data := range machines {
		data := data
		tmp := domain.StockResponse{
			ID:        &data.ID,
			Stock:     data.Stock,
			CreatedBy: data.CreatedBy,
			UpdatedBy: data.UpdatedBy,
			DeletedBy: data.DeletedBy,
			CreatedAt: &data.CreatedAt,
			UpdatedAt: &data.UpdatedAt,
			DeletedAt: &data.DeletedAt,
		}
		result.Stocks = append(result.Stocks, tmp)
	}
	return &result, nil
}
