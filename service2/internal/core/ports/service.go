package ports

import (
	"taobin-service/internal/core/domain"
)

type Service interface {
	UpdateStock(payload domain.StockRequest) (*domain.StockResponse, error)
	GetStock(payload domain.QueryStockRequest) (*domain.StockListResult, error)
}
