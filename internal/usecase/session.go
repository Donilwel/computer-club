package usecase

import (
	"computer-club/internal/repository"
	"computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"log"
	"time"
)

type SessionService interface {
	StartSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models.Session, error)
	EndSession(ctx context.Context, sessionID int64) error
	GetActiveSessions(ctx context.Context) []*models.Session
	MonitorSessions(ctx context.Context)
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

func (u *SessionUsecase) StartSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models.Session, error) {
	_, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	exists, err := u.sessionRepository.HasActiveSession(ctx, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrSessionActive
	}

	pcExists, err := u.computerRepo.IsComputerAvailable(ctx, pcNumber)
	if err != nil {
		return nil, err
	}
	if !pcExists {
		return nil, errors.ErrComputerNotFound
	}

	tariff, err := u.tariffRepo.GetTariffByID(ctx, tariffID)
	if err != nil {
		return nil, errors.ErrTariffNotFound
	}

	balance, err := u.walletRepo.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}
	if balance < tariff.Price {
		return nil, errors.ErrInsufficientFunds
	}

	err = u.walletRepo.Withdraw(nil, userID, tariff.Price)
	if err != nil {
		return nil, err
	}

	_, err = u.walletRepo.CreateTransaction(nil, userID, tariff.Price, string(models.Buy), tariff)
	if err != nil {
		return nil, err
	}

	session, err := u.sessionRepository.CreateSession(ctx, userID, pcNumber, tariffID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *SessionUsecase) EndSession(ctx context.Context, sessionID int64) error {
	return u.sessionRepository.EndSession(ctx, sessionID)
}

func (u *SessionUsecase) GetActiveSessions(ctx context.Context) []*models.Session {
	return u.sessionRepository.GetActiveSessions(ctx)
}

func (u *SessionUsecase) MonitorSessions(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Остановка мониторинга сессий")
			return
		case <-ticker.C:
			u.checkAndCloseExpiredSessions(ctx)
		}
	}
}

func (u *SessionUsecase) checkAndCloseExpiredSessions(ctx context.Context) {
	sessions := u.sessionRepository.GetActiveSessions(ctx)
	now := time.Now()

	for _, session := range sessions {
		if session.EndTime.Before(now) {
			log.Printf("Завершаем сессию %d (пользователь %d)", session.ID, session.UserID)

			// Обновление статуса компьютера
			if err := u.computerRepo.UpdateStatus(ctx, session.PCNumber, models.Free); err != nil {
				log.Printf("Не удалось обновить статус компьютера для сессии %d: %v", session.ID, err)
			}

			// Завершаем сессию
			u.EndSession(ctx, session.ID)
		}
	}
}
