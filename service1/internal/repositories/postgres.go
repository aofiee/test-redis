package repositories

import (
	"errors"
	"taobin-service/internal/core/domain"
	"time"

	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

const (
	layoutDate              = "2006-01-02"
	MachineRefillSubscriber = "MachineRefillSubscriber"
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

func (r *Postgres) CreateMachine(payload domain.MachineRequest) (*domain.MachineResponse, error) {
	var (
		response domain.MachineResponse
	)
	machine := domain.Machine{
		Machine:   payload.Machine,
		Stock:     payload.Stock,
		CreatedBy: payload.CreatedBy,
		UpdatedBy: payload.UpdatedBy,
	}
	err := r.dbGorm.Create(&machine).Error
	if err != nil {
		return &domain.MachineResponse{}, err
	}
	response = domain.MachineResponse{
		ID:        &machine.ID,
		Machine:   machine.Machine,
		Stock:     machine.Stock,
		CreatedBy: machine.CreatedBy,
		UpdatedBy: machine.UpdatedBy,
		DeletedBy: machine.DeletedBy,
		CreatedAt: &machine.CreatedAt,
		UpdatedAt: &machine.UpdatedAt,
		DeletedAt: &machine.DeletedAt,
	}
	return &response, nil
}

func (r *Postgres) UpdateMachine(payload domain.MachineRequest) (*domain.MachineResponse, error) {
	var (
		response domain.MachineResponse
		result   domain.Machine
	)
	query := domain.QueryMachineRequest{
		ID: payload.ID,
	}
	condition := r.machineCondition(query)
	columns := r.updateMachineColumns(payload)
	tx := r.dbGorm.Begin()
	tx.Table(domain.Machine{}.TableName()).Where(condition).Updates(columns)
	if tx.Error != nil {
		tx.Rollback()
		return &domain.MachineResponse{}, tx.Error
	}
	tx.Where(condition).First(&result)
	if tx.Error != nil {
		tx.Rollback()
		return &domain.MachineResponse{}, tx.Error
	}
	tx.Commit()
	if result.ID == 0 {
		return &domain.MachineResponse{}, errors.New("Not found")
	}
	response = domain.MachineResponse{
		ID:        &result.ID,
		Machine:   result.Machine,
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

func (r *Postgres) machineCondition(condition domain.QueryMachineRequest) map[string]interface{} {
	expression := map[string]interface{}{}
	if condition.ID != nil {
		expression["id"] = condition.ID
	}
	if condition.Machine != nil {
		expression["machine"] = condition.Machine
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

func (r *Postgres) updateMachineColumns(payload domain.MachineRequest) map[string]interface{} {
	expression := map[string]interface{}{}
	if payload.Machine != nil {
		expression["machine"] = payload.Machine
	}
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

func (r *Postgres) DeleteMachine(payload domain.MachineRequest) (*domain.MachineResponse, error) {
	var (
		response domain.MachineResponse
		result   domain.Machine
	)
	query := domain.QueryMachineRequest{
		ID: payload.ID,
	}
	condition := r.machineCondition(query)
	columns := r.updateMachineColumns(payload)
	if len(columns) == 0 {
		return &domain.MachineResponse{}, errors.New("Not update columns")
	}
	if err := r.dbGorm.Where(condition).First(&result).Error; err != nil {
		return &domain.MachineResponse{}, err
	}
	tx := r.dbGorm.Begin()
	tx.Model(&result).Updates(columns)
	if tx.Error != nil {
		tx.Rollback()
		return &domain.MachineResponse{}, tx.Error
	}
	tx.Delete(&result)
	if tx.Error != nil {
		tx.Rollback()
		return &domain.MachineResponse{}, tx.Error
	}
	tx.Commit()
	response = domain.MachineResponse{
		ID:        &result.ID,
		Machine:   result.Machine,
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

func (r *Postgres) GetMachine(payload domain.QueryMachineRequest) (*domain.MachineListResult, error) {
	var (
		machines []domain.Machine
	)
	condition := r.machineCondition(payload)
	tx := r.dbGorm
	count := r.dbGorm
	if payload.DeletedBy != nil || payload.DeletedAt != nil {
		tx = tx.Unscoped().Where(condition)
		count = count.Unscoped().Select("COUNT(id)").Where(condition)
	} else {
		tx = tx.Where(condition)
		count = count.Model(&domain.Machine{}).Select("COUNT(id)").Where(condition)
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
	result := domain.MachineListResult{
		Machines: []domain.MachineResponse{},
	}
	result.CurrentPage = payload.Page
	result.PerPage = &payload.Pagination.Limit
	var toatalItem int64
	count.Count(&toatalItem)
	result.TotalItem = &toatalItem
	for _, data := range machines {
		data := data
		tmp := domain.MachineResponse{
			ID:        &data.ID,
			Machine:   data.Machine,
			Stock:     data.Stock,
			CreatedBy: data.CreatedBy,
			UpdatedBy: data.UpdatedBy,
			DeletedBy: data.DeletedBy,
			CreatedAt: &data.CreatedAt,
			UpdatedAt: &data.UpdatedAt,
			DeletedAt: &data.DeletedAt,
		}
		result.Machines = append(result.Machines, tmp)
	}
	return &result, nil
}
