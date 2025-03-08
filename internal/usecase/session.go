package usecase

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"computer-club/internal/repository"
	"context"
	"log"
	"time"
)

type SessionService interface {
	StartSession(userID int64, pcNumber int, tariffID int64) (*models.Session, error)
	EndSession(sessionID int64) error
	GetActiveSessions() []*models.Session
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
	walletService WalletService) *SessionUsecase {
	return &SessionUsecase{sessionRepository: sessionRepository,
		userRepo:      userRepo,
		computerRepo:  computerRepo,
		walletService: walletService}
}

func (u *SessionUsecase) StartSession(userID int64, pcNumber int, tariffID int64) (*models.Session, error) {
	_, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	session, err := u.sessionRepository.StartSession(userID, pcNumber, tariffID)
	if err != nil {
		return nil, err
	}

	err = u.walletService.ChargeForSession(userID, tariffID)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *SessionUsecase) EndSession(sessionID int64) error {
	return u.sessionRepository.EndSession(sessionID)
}

func (u *SessionUsecase) GetActiveSessions() []*models.Session {
	return u.sessionRepository.GetActiveSessions()
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
			u.checkAndCloseExpiredSessions()
		}
	}
}

func (u *SessionUsecase) checkAndCloseExpiredSessions() {
	sessions := u.sessionRepository.GetActiveSessions()
	now := time.Now()

	for _, session := range sessions {
		if session.EndTime.Before(now) {
			log.Printf("Завершаем сессию %d (пользователь %d)", session.ID, session.UserID)

			// Обновление статуса компьютера
			if err := u.computerRepo.UpdateStatus(session.PCNumber, models.Free); err != nil {
				log.Printf("Не удалось обновить статус компьютера для сессии %d: %v", session.ID, err)
			}

			// Завершаем сессию
			u.EndSession(session.ID)
		}
	}
}
