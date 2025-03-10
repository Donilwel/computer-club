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

	tx := u.sessionRepository.BeginTransaction(ctx)
	defer tx.Rollback()

	exists, err := u.sessionRepository.HasActiveSession(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.ErrSessionActive
	}

	pcBusy, err := u.computerRepo.IsComputerAvailable(ctx, tx, pcNumber)
	if err != nil {
		return nil, err
	}
	if pcBusy {
		return nil, errors.ErrPCBusy
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

	err = u.walletRepo.Withdraw(ctx, tx, userID, tariff.Price)
	if err != nil {
		return nil, err
	}

	_, err = u.walletRepo.CreateTransaction(ctx, tx, userID, tariff.Price, string(models.Buy), tariff)
	if err != nil {
		return nil, err
	}

	session, err := u.sessionRepository.CreateSession(ctx, tx, userID, pcNumber, tariffID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.ErrCommitData
	}

	err = u.sessionRepository.CacheSession(ctx, session)
	if err != nil {
		return nil, errors.ErrCacheSession
	}

	return session, nil
}

func (u *SessionUsecase) EndSession(ctx context.Context, sessionID int64) error {
	tx := u.sessionRepository.BeginTransaction(ctx)
	defer tx.Rollback()

	session, err := u.sessionRepository.GetSessionByID(ctx, tx, sessionID)
	if err != nil {
		return errors.ErrSessionNotFound
	}

	if session.Status == models.Finished {
		return errors.ErrStatusSessionAlreadyFinished
	}

	err = u.sessionRepository.MarkSessionFinished(ctx, tx, sessionID)
	if err != nil {
		return errors.ErrUpdateSession
	}

	err = u.computerRepo.MarkComputerFree(ctx, tx, session.PCNumber)
	if err != nil {
		return errors.ErrUpdateComputer
	}

	err = tx.Commit()
	if err != nil {
		return errors.ErrCommitData
	}

	err = u.sessionRepository.DeleteSessionCache(ctx, sessionID)
	if err != nil {
		return errors.ErrDeleteRedis
	}

	return nil
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
			if err := u.computerRepo.MarkComputerFree(ctx, nil, session.PCNumber); err != nil {
				log.Printf("Не удалось обновить статус компьютера для сессии %d: %v", session.ID, err)
			}

			// Завершаем сессию
			u.EndSession(ctx, session.ID)
		}
	}
}
