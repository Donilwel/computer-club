package di

import (
	"computer-club/config"
	"computer-club/internal/delivery/httpService"
	"computer-club/internal/logger"
	"computer-club/internal/middleware"
	"computer-club/internal/repository"
	"computer-club/internal/usecase"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

// Container отвечает за хранение и инициализацию зависимостей
type Container struct {
	Cfg             *config.Config
	Log             *logrus.Logger
	DB              *gorm.DB
	RedisClient     *redis.Client
	UserRepo        repository.UserRepository
	SessionRepo     repository.SessionRepository
	ComputerRepo    repository.ComputerRepository
	TariffRepo      repository.TariffRepository
	WalletRepo      repository.WalletRepository
	UserUsecase     *usecase.UserUsecase
	SessionUsecase  *usecase.SessionUsecase
	ComputerUsecase *usecase.ComputerUsecase
	TariffUsecase   *usecase.TariffUsecase
	WalletUsecase   *usecase.WalletUsecase
	Router          *chi.Mux
}

// NewContainer создает новый контейнер зависимостей
func NewContainer() *Container {
	cfg := config.LoadConfig()
	log := logger.NewLogger()

	// Подключение к БД и Redis
	db := repository.NewPostgresDB(cfg)
	redisClient := repository.NewRedisClient(cfg)
	repository.Migrate(db)

	// Инициализация репозиториев
	userRepo := repository.NewPostgresUserRepo(db)
	sessionRepo := repository.NewPostgresSessionRepo(db, redisClient)
	computerRepo := repository.NewComputerRepository(db)
	tariffRepo := repository.NewTariffRepositoryPostgres(db)
	walletRepo := repository.NewPostgresWalletRepo(db)

	// Инициализация usecase'ов
	tariffUsecase := usecase.NewTariffUsecase(tariffRepo)
	walletUsecase := usecase.NewWalletUsecase(walletRepo, tariffUsecase, userRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, walletUsecase)
	sessionUsecase := usecase.NewSessionUsecase(sessionRepo, userRepo, computerRepo, walletUsecase)
	computerUsecase := usecase.NewComputerUsecase(computerRepo)

	// Инициализация HTTP-хендлера
	handler := httpService.NewHandler(userUsecase, computerUsecase, sessionUsecase, tariffUsecase, walletUsecase, log)
	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware(log))
	handler.RegisterRoutes(r)

	return &Container{
		Cfg:             cfg,
		Log:             log,
		DB:              db,
		RedisClient:     redisClient,
		UserRepo:        userRepo,
		SessionRepo:     sessionRepo,
		ComputerRepo:    computerRepo,
		TariffRepo:      tariffRepo,
		WalletRepo:      walletRepo,
		UserUsecase:     userUsecase,
		SessionUsecase:  sessionUsecase,
		ComputerUsecase: computerUsecase,
		TariffUsecase:   tariffUsecase,
		WalletUsecase:   walletUsecase,
		Router:          r,
	}
}

// RunServer запускает HTTP-сервер
func (c *Container) RunServer(ctx context.Context) {
	// Запускаем мониторинг сессий
	go c.SessionUsecase.MonitorSessions(ctx)

	fmt.Println("Server started on :", c.Cfg.ServerPort)
	err := http.ListenAndServe(":"+c.Cfg.ServerPort, c.Router)
	if err != nil {
		c.Log.Fatal("Server error: ", err)
	}
}
