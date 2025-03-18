package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"computer-club/pkg/errors"
	"context"
)

type ComputerService interface {
	GetComputersStatus(ctx context.Context) ([]*models.Computer, error)
	DeleteComputer(ctx context.Context, id int) error
	AddComputer(ctx context.Context) (*models.Computer, error)
}

type ComputerUsecase struct {
	computerRepo repository.ComputerRepository
}

func NewComputerUsecase(computerRepo repository.ComputerRepository) ComputerService {
	return &ComputerUsecase{computerRepo: computerRepo}
}

func (u *ComputerUsecase) DeleteComputer(ctx context.Context, id int) error {
	computer, err := u.computerRepo.GetComputerByID(ctx, id)
	if err != nil {
		return err
	}

	if computer.Status == models.Busy {
		return errors.ErrPCBusy
	}

	return u.computerRepo.DeleteComputer(ctx, computer)
}
func (u *ComputerUsecase) GetComputersStatus(ctx context.Context) ([]*models.Computer, error) {
	return u.computerRepo.GetComputers(ctx)
}

func (u *ComputerUsecase) AddComputer(ctx context.Context) (*models.Computer, error) {
	return u.computerRepo.AddComputer(ctx)
}
