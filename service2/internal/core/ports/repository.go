package ports

import (
	"taobin-service/internal/core/domain"
)

type Repository interface {
	UpdateStock(payload domain.StockRequest) (*domain.StockResponse, error)
	GetStock(payload domain.QueryStockRequest) (*domain.StockListResult, error)
}
