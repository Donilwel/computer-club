package repository

import (
	"computer-club/internal/repository/models"
	"computer-club/pkg/errors"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sync"
	"time"
)

type SessionRepository interface {
	EndSession(ctx context.Context, sessionID int64) error
	GetActiveSessions(ctx context.Context) []*models.Session
	GetSessionByID(ctx context.Context, sessionID int64) (*models.Session, error)
	CheckStatus(session models.Session, status string) error
	HasActiveSession(ctx context.Context, userID int64) (bool, error)
	CreateSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models.Session, error)
}

type PostgresSessionRepo struct {
	db    *gorm.DB
	redis *redis.Client
	mu    sync.Mutex
}

func NewPostgresSessionRepo(db *gorm.DB, redis *redis.Client) SessionRepository {
	return &PostgresSessionRepo{db: db, redis: redis}
}

func (r *PostgresSessionRepo) GetSessionByID(ctx context.Context, sessionID int64) (*models.Session, error) {
	var session models.Session
	if err := r.db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error; err != nil {
		return nil, errors.ErrSessionNotFound
	}
	return &session, nil
}

func (r *PostgresSessionRepo) HasActiveSession(ctx context.Context, userID int64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Session{}).Where("user_id = ? AND status = ?", userID, models.Active).Count(&count).Error
	return count > 0, err
}

func (r *PostgresSessionRepo) CheckStatus(sesion models.Session, status string) error {
	if sesion.Status != models.SessionStatus(status) {
		return errors.ErrFailedStatus
	}
	return nil
}

func (r *PostgresSessionRepo) EndSession(ctx context.Context, sessionID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return errors.ErrStartTransaction
	}

	// Находим сессию
	var session models.Session
	if err := tx.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error; err != nil {
		tx.Rollback()
		return errors.ErrSessionNotFound
	}

	// Завершаем сессию
	if err := tx.WithContext(ctx).Model(&models.Session{}).Where("id = ?", sessionID).Update("status", models.Finished).Error; err != nil {
		tx.Rollback()
		return errors.ErrUpdateSession
	}

	// Освобождаем компьютер
	if err := tx.WithContext(ctx).Model(&models.Computer{}).Where("pc_number = ?", session.PCNumber).Update("status", models.Free).Error; err != nil {
		tx.Rollback()
		return errors.ErrUpdateComputer
	}

	// Удаляем сессию из кеша
	if err := r.redis.Del(ctx, getSessionKey(sessionID)).Err(); err != nil {
		tx.Rollback()
		return errors.ErrDeleteRedis
	}
	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return errors.ErrCommitData
	}

	return nil
}

func (r *PostgresSessionRepo) GetActiveSessions(ctx context.Context) []*models.Session {
	// Проверяем кеш
	var sessions []*models.Session
	keys, _ := r.redis.Keys(ctx, "session:*").Result()
	if len(keys) > 0 {
		for _, key := range keys {
			var session models.Session
			sessionJSON, _ := r.redis.Get(ctx, key).Result()
			json.Unmarshal([]byte(sessionJSON), &session)
			sessions = append(sessions, &session)
		}
		return sessions
	}

	// Если в кеше нет, загружаем из БД
	r.db.WithContext(ctx).Where("status = ?", models.Active).Find(&sessions)

	// Кешируем результат
	for _, session := range sessions {
		sessionJSON, _ := json.Marshal(session)
		r.redis.Set(ctx, getSessionKey(session.ID), sessionJSON, 10*time.Minute)
	}

	return sessions
}

func (r *PostgresSessionRepo) CreateSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models.Session, error) {
	startTime := time.Now()
	endTime := startTime.Add(2 * time.Hour)

	session := &models.Session{
		UserID:    userID,
		PCNumber:  pcNumber,
		TariffID:  tariffID,
		Status:    models.Active,
		StartTime: startTime,
		EndTime:   &endTime,
	}

	if err := r.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, errors.ErrCreatedSession
	}

	// Обновляем статус ПК
	if err := r.db.WithContext(ctx).Model(&models.Computer{}).
		Where("pc_number = ?", pcNumber).
		Update("status", models.Busy).Error; err != nil {
		return nil, errors.ErrUpdateComputerStatus
	}

	return session, nil
}

// Вспомогательная функция для генерации ключа Redis
func getSessionKey(sessionID int64) string {
	return "session:" + fmt.Sprint(sessionID)
}
