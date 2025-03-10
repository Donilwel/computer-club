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
	GetActiveSessions(ctx context.Context) []*models.Session
	GetSessionByID(ctx context.Context, tx Transaction, sessionID int64) (*models.Session, error)
	CheckStatus(session models.Session, status string) error
	HasActiveSession(ctx context.Context, tx Transaction, userID int64) (bool, error)
	CreateSession(ctx context.Context, tx Transaction, userID int64, pcNumber int, tariffID int64) (*models.Session, error)
	MarkSessionFinished(ctx context.Context, tx Transaction, sessionID int64) error
	CacheSession(ctx context.Context, session *models.Session) error
	DeleteSessionCache(ctx context.Context, id int64) error
	BeginTransaction(ctx context.Context) Transaction
}

type PostgresSessionRepo struct {
	db    *gorm.DB
	redis *redis.Client
	mu    sync.Mutex
}

func NewPostgresSessionRepo(db *gorm.DB, redis *redis.Client) SessionRepository {
	return &PostgresSessionRepo{db: db, redis: redis}
}

func (r *PostgresSessionRepo) BeginTransaction(ctx context.Context) Transaction {
	return &GormTransaction{tx: r.db.WithContext(ctx).Begin()}
}

func (r *PostgresSessionRepo) GetSessionByID(ctx context.Context, tx Transaction, sessionID int64) (*models.Session, error) {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	var session models.Session
	if err := db.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error; err != nil {
		return nil, errors.ErrSessionNotFound
	}
	return &session, nil
}

func (r *PostgresSessionRepo) HasActiveSession(ctx context.Context, tx Transaction, userID int64) (bool, error) {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	var count int64
	err := db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND status = ?", userID, models.Active).
		Count(&count).Error
	return count > 0, err
}

func (r *PostgresSessionRepo) CheckStatus(session models.Session, status string) error {
	if session.Status != models.SessionStatus(status) {
		return errors.ErrFailedStatus
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

func (r *PostgresSessionRepo) CreateSession(ctx context.Context, tx Transaction, userID int64, pcNumber int, tariffID int64) (*models.Session, error) {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
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

	if err := db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, errors.ErrCreatedSession
	}

	// Обновляем статус ПК
	if err := db.WithContext(ctx).Model(&models.Computer{}).
		Where("pc_number = ?", pcNumber).
		Update("status", models.Busy).Error; err != nil {
		return nil, errors.ErrUpdateComputerStatus
	}

	return session, nil
}

func (r *PostgresSessionRepo) MarkSessionFinished(ctx context.Context, tx Transaction, sessionID int64) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	err := db.WithContext(ctx).Model(&models.Session{}).
		Where("id = ?", sessionID).
		Update("status", models.Finished).Error
	if err != nil {
		return errors.ErrUpdateSession
	}
	return nil
}
func (r *PostgresSessionRepo) CacheSession(ctx context.Context, session *models.Session) error {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}
	cacheKey := getSessionKey(session.ID)

	err = r.redis.Set(ctx, cacheKey, sessionJSON, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresSessionRepo) DeleteSessionCache(ctx context.Context, sessionID int64) error {
	return r.redis.Del(ctx, getSessionKey(sessionID)).Err()
}

func getSessionKey(sessionID int64) string {
	return "session:" + fmt.Sprint(sessionID)
}
