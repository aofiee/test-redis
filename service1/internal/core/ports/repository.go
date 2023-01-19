package ports

import (
	"taobin-service/internal/core/domain"
)

type Repository interface {
	CreateMachine(payload domain.MachineRequest) (*domain.MachineResponse, error)
	UpdateMachine(payload domain.MachineRequest) (*domain.MachineResponse, error)
	DeleteMachine(payload domain.MachineRequest) (*domain.MachineResponse, error)
	GetMachine(payload domain.QueryMachineRequest) (*domain.MachineListResult, error)
}
