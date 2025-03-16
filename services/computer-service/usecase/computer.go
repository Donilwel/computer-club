package usecase

import (
	"computer-club/internal/repository"
	"computer-club/internal/repository/models"
	"context"
)

type ComputerService interface {
	GetComputersStatus(ctx context.Context) ([]models.Computer, error)
}

type ComputerUsecase struct {
	computerRepo repository.ComputerRepository
}

func NewComputerUsecase(computerRepo repository.ComputerRepository) ComputerService {
	return &ComputerUsecase{computerRepo: computerRepo}
}

func (u *ComputerUsecase) GetComputersStatus(ctx context.Context) ([]models.Computer, error) {
	return u.computerRepo.GetComputers(ctx)
}
