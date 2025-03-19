package usecase

import (
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"computer-club/pkg/errors"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type SessionService interface {
	StartSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models.Session, error)
	EndSession(ctx context.Context, sessionID int64) error
	GetActiveSessions(ctx context.Context) []*models.Session
	CheckExpiredSessions(ctx context.Context)
}

type SessionUsecase struct {
	sessionRepository repository.SessionRepository
	userRepo          repository.UserRepository
	computerRepo      repository.ComputerRepository
	tariffRepo        repository.TariffRepository
	walletRepo        repository.WalletRepository
}

func NewSessionUsecase(sessionRepository repository.SessionRepository,
	userRepo repository.UserRepository,
	computerRepo repository.ComputerRepository,
	tariffRepo repository.TariffRepository,
	walletRepo repository.WalletRepository) SessionService {
	return &SessionUsecase{sessionRepository: sessionRepository,
		userRepo:     userRepo,
		computerRepo: computerRepo,
		tariffRepo:   tariffRepo,
		walletRepo:   walletRepo}
}

func (u *SessionUsecase) StartSession(ctx context.Context, userID int64,
	pcNumber int, tariffID int64) (*models.Session, error) {
	_, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	tx := u.sessionRepository.BeginTransaction(ctx)
	defer tx.Rollback()

	exists, err := u.sessionRepository.HasActiveSession(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.ErrSessionActive
	}

	computer, err := u.computerRepo.GetComputerByID(ctx, pcNumber)
	if err != nil {
		return nil, err
	}

	if computer.Status == models.Busy {
		return nil, errors.ErrPCBusy
	}

	tariff, err := u.tariffRepo.GetTariffByID(ctx, tariffID)
	if err != nil {
		return nil, err
	}

	balance, err := u.walletRepo.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}
	if balance < tariff.Price {
		return nil, errors.ErrInsufficientFunds
	}

	err = u.walletRepo.Withdraw(ctx, tx, userID, tariff.Price)
	if err != nil {
		return nil, err
	}

	_, err = u.walletRepo.CreateTransaction(ctx, tx, userID, tariff.Price, string(models.Buy), tariff)
	if err != nil {
		return nil, err
	}

	session, err := u.sessionRepository.CreateSession(ctx, tx, userID, pcNumber, tariff)
	if err != nil {
		return nil, err
	}

	err = u.computerRepo.ChangeComputerStatus(ctx, tx, computer, string(models.Busy))
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	err = u.sessionRepository.CacheSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *SessionUsecase) EndSession(ctx context.Context, sessionID int64) error {
	tx := u.sessionRepository.BeginTransaction(ctx)
	defer tx.Rollback()

	session, err := u.sessionRepository.GetSessionByID(ctx, tx, sessionID)
	if err != nil {
		return err
	}

	if session.Status == models.Finished {
		return errors.ErrStatusSessionAlreadyFinished
	}

	err = u.sessionRepository.MarkSessionFinished(ctx, tx, sessionID)
	if err != nil {
		return err
	}

	computer, err := u.computerRepo.GetComputerByID(ctx, session.PCNumber)
	if err != nil {
		return err
	}

	if computer.Status == models.Free {
		return errors.ErrComputerAlreadyFree
	}

	err = u.computerRepo.ChangeComputerStatus(ctx, tx, computer, string(models.Free))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	err = u.sessionRepository.DeleteSessionCache(ctx, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionUsecase) GetActiveSessions(ctx context.Context) []*models.Session {
	return u.sessionRepository.GetActiveSessions(ctx)
}

func (u *SessionUsecase) CheckExpiredSessions(ctx context.Context) {
	sessions := u.sessionRepository.GetActiveSessions(ctx)
	now := time.Now()

	for _, session := range sessions {
		if session.EndTime != nil && session.EndTime.Before(now) {
			logrus.Print(fmt.Sprintf("Сессия %d просрочена, завершаем...\n", session.ID))
			err := u.EndSession(ctx, session.ID)
			if err != nil {
				logrus.Print(fmt.Sprintf("Ошибка при завершении сессии %d: %v\n", session.ID, err))
			}
		}
	}
}
