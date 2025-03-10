package repository

import (
	models2 "computer-club/internal/repository/models"
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
	StartSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models2.Session, error)
	EndSession(ctx context.Context, sessionID int64) error
	GetActiveSessions(ctx context.Context) []*models2.Session
}

type PostgresSessionRepo struct {
	db    *gorm.DB
	redis *redis.Client
	mu    sync.Mutex
}

func NewPostgresSessionRepo(db *gorm.DB, redis *redis.Client) SessionRepository {
	return &PostgresSessionRepo{db: db, redis: redis}
}

func (r *PostgresSessionRepo) StartSession(ctx context.Context, userID int64, pcNumber int, tariffID int64) (*models2.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, errors.ErrStartTransaction
	}

	// Проверяем, есть ли активная сессия у пользователя
	var activeSession models2.Session
	if err := tx.WithContext(ctx).Where("user_id = ? AND status = ?", userID, models2.Active).First(&activeSession).Error; err == nil {
		tx.Rollback()
		return nil, errors.ErrSessionActive
	}

	var tariff models2.Tariff
	if err := tx.WithContext(ctx).First(&tariff, tariffID).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrTariffNotFound
	}

	// Проверяем, существует ли этот ПК
	var computer models2.Computer
	if err := tx.WithContext(ctx).Where("pc_number = ?", pcNumber).First(&computer).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrComputerNotFound
	}

	// Проверяем, свободен ли компьютер
	if computer.Status == models2.Busy {
		tx.Rollback()
		return nil, errors.ErrPCBusy
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(tariff.Duration) * time.Minute)

	// Создаем новую сессию
	session := &models2.Session{
		UserID:    userID,
		PCNumber:  pcNumber,
		TariffID:  tariffID,
		Status:    models2.Active,
		StartTime: startTime,
		EndTime:   &endTime,
	}
	if err := tx.Create(session).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrCreatedSession
	}

	// Обновляем статус компьютера
	if err := tx.WithContext(ctx).Model(&models2.Computer{}).Where("pc_number = ?", pcNumber).Update("status", models2.Busy).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrUpdateComputerStatus
	}

	// Кешируем активную сессию в Redis
	sessionJSON, _ := json.Marshal(session)
	r.redis.Set(ctx, getSessionKey(session.ID), sessionJSON, 10*time.Minute)

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return nil, errors.ErrCommitData
	}

	return session, nil
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
	var session models2.Session
	if err := tx.WithContext(ctx).Where("id = ?", sessionID).First(&session).Error; err != nil {
		tx.Rollback()
		return errors.ErrSessionNotFound
	}

	// Завершаем сессию
	if err := tx.WithContext(ctx).Model(&models2.Session{}).Where("id = ?", sessionID).Update("status", models2.Finished).Error; err != nil {
		tx.Rollback()
		return errors.ErrUpdateSession
	}

	// Освобождаем компьютер
	if err := tx.WithContext(ctx).Model(&models2.Computer{}).Where("pc_number = ?", session.PCNumber).Update("status", models2.Free).Error; err != nil {
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

func (r *PostgresSessionRepo) GetActiveSessions(ctx context.Context) []*models2.Session {
	// Проверяем кеш
	var sessions []*models2.Session
	keys, _ := r.redis.Keys(ctx, "session:*").Result()
	if len(keys) > 0 {
		for _, key := range keys {
			var session models2.Session
			sessionJSON, _ := r.redis.Get(ctx, key).Result()
			json.Unmarshal([]byte(sessionJSON), &session)
			sessions = append(sessions, &session)
		}
		return sessions
	}

	// Если в кеше нет, загружаем из БД
	r.db.WithContext(ctx).Where("end_time IS NULL").Find(&sessions)

	// Кешируем результат
	for _, session := range sessions {
		sessionJSON, _ := json.Marshal(session)
		r.redis.Set(ctx, getSessionKey(session.ID), sessionJSON, 10*time.Minute)
	}

	return sessions
}

// Вспомогательная функция для генерации ключа Redis
func getSessionKey(sessionID int64) string {
	return "session:" + fmt.Sprint(sessionID)
}
