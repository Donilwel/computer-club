package repository

import (
	"computer-club/internal/errors"
	"computer-club/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sync"
	"time"
)

type SessionRepository interface {
	StartSession(userID int64, pcNumber int, tariffID int64) (*models.Session, error)
	EndSession(sessionID int64) error
	GetActiveSessions() []*models.Session
}

type PostgresSessionRepo struct {
	db    *gorm.DB
	redis *redis.Client
	mu    sync.Mutex
}

func NewPostgresSessionRepo(db *gorm.DB, redis *redis.Client) SessionRepository {
	return &PostgresSessionRepo{db: db, redis: redis}
}

func (r *PostgresSessionRepo) StartSession(userID int64, pcNumber int, tariffID int64) (*models.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, errors.ErrStartTransaction
	}

	// Проверяем, есть ли активная сессия у пользователя
	var activeSession models.Session
	if err := tx.Where("user_id = ? AND status = ?", userID, models.Active).First(&activeSession).Error; err == nil {
		tx.Rollback()
		return nil, errors.ErrSessionActive
	}

	var tariff models.Tariff
	if err := tx.First(&tariff, tariffID).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrTariffNotFound
	}

	// Проверяем, существует ли этот ПК
	var computer models.Computer
	if err := tx.Where("pc_number = ?", pcNumber).First(&computer).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrComputerNotFound
	}

	// Проверяем, свободен ли компьютер
	if computer.Status == models.Busy {
		tx.Rollback()
		return nil, errors.ErrPCBusy
	}

	startTime := time.Now()
	endTime := startTime.Add(time.Duration(tariff.Duration) * time.Minute)

	// Создаем новую сессию
	session := &models.Session{
		UserID:    userID,
		PCNumber:  pcNumber,
		TariffID:  tariffID,
		Status:    models.Active,
		StartTime: startTime,
		EndTime:   &endTime,
	}
	if err := tx.Create(session).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrCreatedSession
	}

	// Обновляем статус компьютера
	if err := tx.Model(&models.Computer{}).Where("pc_number = ?", pcNumber).Update("status", models.Busy).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrUpdateComputerStatus
	}

	// Кешируем активную сессию в Redis
	ctx := context.Background()
	sessionJSON, _ := json.Marshal(session)
	r.redis.Set(ctx, getSessionKey(session.ID), sessionJSON, 10*time.Minute)

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		return nil, errors.ErrCommitData
	}

	return session, nil
}

func (r *PostgresSessionRepo) EndSession(sessionID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Начинаем транзакцию
	tx := r.db.Begin()
	if tx.Error != nil {
		return errors.ErrStartTransaction
	}

	// Находим сессию
	var session models.Session
	if err := tx.Where("id = ?", sessionID).First(&session).Error; err != nil {
		tx.Rollback()
		return errors.ErrSessionNotFound
	}

	// Завершаем сессию
	if err := tx.Model(&models.Session{}).Where("id = ?", sessionID).Update("status", models.Finished).Error; err != nil {
		tx.Rollback()
		return errors.ErrUpdateSession
	}

	// Освобождаем компьютер
	if err := tx.Model(&models.Computer{}).Where("pc_number = ?", session.PCNumber).Update("status", models.Free).Error; err != nil {
		tx.Rollback()
		return errors.ErrUpdateComputer
	}

	// Удаляем сессию из кеша
	ctx := context.Background()
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

func (r *PostgresSessionRepo) GetActiveSessions() []*models.Session {
	ctx := context.Background()

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
	r.db.Where("end_time IS NULL").Find(&sessions)

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
