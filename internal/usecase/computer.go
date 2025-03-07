package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
)

type ComputerService interface {
	GetComputersStatus() ([]models.Computer, error)
}

type ComputerUsecase struct {
	computerRepo repository.ComputerRepository
}

func NewComputerUsecase(computerRepo repository.ComputerRepository) *ComputerUsecase {
	return &ComputerUsecase{computerRepo: computerRepo}
}

func (u *ComputerUsecase) GetComputersStatus() ([]models.Computer, error) {
	return u.computerRepo.GetComputers()
}
