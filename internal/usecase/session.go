package usecase

import (
	"computer-club/internal/repository"
	models2 "computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"log"
	"time"
)

type SessionService interface {
	StartSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models2.Session, error)
	EndSession(ctx context.Context, sessionID int64) error
	GetActiveSessions(ctx context.Context) []*models2.Session
	MonitorSessions(ctx context.Context)
}

type SessionUsecase struct {
	sessionRepository repository.SessionRepository
	userRepo          repository.UserRepository
	computerRepo      repository.ComputerRepository
	walletService     WalletService
}

func NewSessionUsecase(sessionRepository repository.SessionRepository,
	userRepo repository.UserRepository,
	computerRepo repository.ComputerRepository,
	walletService WalletService) SessionService {
	return &SessionUsecase{sessionRepository: sessionRepository,
		userRepo:      userRepo,
		computerRepo:  computerRepo,
		walletService: walletService}
}

func (u *SessionUsecase) StartSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models2.Session, error) {
	_, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	session, err := u.sessionRepository.StartSession(ctx, userID, pcNumber, tariffID)
	if err != nil {
		return nil, err
	}

	err = u.walletService.ChargeForSession(ctx, userID, tariffID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *SessionUsecase) EndSession(ctx context.Context, sessionID int64) error {
	return u.sessionRepository.EndSession(ctx, sessionID)
}

func (u *SessionUsecase) GetActiveSessions(ctx context.Context) []*models2.Session {
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
			if err := u.computerRepo.UpdateStatus(ctx, session.PCNumber, models2.Free); err != nil {
				log.Printf("Не удалось обновить статус компьютера для сессии %d: %v", session.ID, err)
			}

			// Завершаем сессию
			u.EndSession(ctx, session.ID)
		}
	}
}
