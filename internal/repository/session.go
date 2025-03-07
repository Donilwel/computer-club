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
	StartSession(userID int64, pcNumber int) (*models.Session, error)
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

func (r *PostgresSessionRepo) StartSession(userID int64, pcNumber int) (*models.Session, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Проверяем, существует ли этот ПК
	var computer models.Computer
	if err := r.db.Where("pc_number = ?", pcNumber).First(&computer).Error; err != nil {
		return nil, errors.ErrComputerNotFound
	}

	// Проверяем, свободен ли компьютер
	if computer.Status == models.Busy {
		return nil, errors.ErrPCBusy
	}

	// Создаем новую сессию
	session := &models.Session{
		UserID:    userID,
		PCNumber:  pcNumber,
		StartTime: time.Now(),
	}
	if err := r.db.Create(session).Error; err != nil {
		return nil, errors.ErrCreatedSession
	}

	// Обновляем статус компьютера
	r.db.Model(&models.Computer{}).Where("pc_number = ?", pcNumber).Update("status", models.Busy)

	// Кешируем активную сессию в Redis
	ctx := context.Background()
	sessionJSON, _ := json.Marshal(session)
	r.redis.Set(ctx, getSessionKey(session.ID), sessionJSON, 10*time.Minute)

	return session, nil
}

func (r *PostgresSessionRepo) EndSession(sessionID int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Находим сессию
	var session models.Session
	if err := r.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return errors.ErrSessionNotFound
	}

	// Завершаем сессию
	now := time.Now()
	if err := r.db.Model(&models.Session{}).Where("id = ?", sessionID).Update("end_time", now).Error; err != nil {
		return errors.ErrUpdateSession
	}

	// Освобождаем компьютер
	if err := r.db.Model(&models.Computer{}).Where("pc_number = ?", session.PCNumber).Update("status", models.Free).Error; err != nil {
		return errors.ErrUpdateComputer
	}

	// Удаляем сессию из кеша
	ctx := context.Background()
	r.redis.Del(ctx, getSessionKey(sessionID))

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
