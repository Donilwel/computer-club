package di

import (
	"computer-club/internal/config"
	"computer-club/internal/delivery/httpService"
	"computer-club/internal/handlers"
	"computer-club/internal/middleware"
	"computer-club/internal/repository"
	"computer-club/internal/usecase"
	"computer-club/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Container отвечает за хранение и инициализацию зависимостей
type Container struct {
	Cfg             *config.Config
	Log             *logrus.Logger
	DB              *gorm.DB
	RedisClient     *redis.Client
	UserHandler     handlers.UserHandler
	SessionHandler  handlers.SessionHandler
	ComputerHandler handlers.ComputerHandler
	TariffHandler   handlers.TariffHandler
	WalletHandler   handlers.WalletHandler
	UserRepo        repository.UserRepository
	SessionRepo     repository.SessionRepository
	ComputerRepo    repository.ComputerRepository
	TariffRepo      repository.TariffRepository
	WalletRepo      repository.WalletRepository
	UserUsecase     *usecase.UserService
	SessionUsecase  *usecase.SessionService
	ComputerUsecase *usecase.ComputerService
	TariffUsecase   *usecase.TariffService
	WalletUsecase   *usecase.WalletService
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

	// Инициализация хендлеров
	userHandler := handlers.NewUserHandler(userUsecase, log) // Добавили log
	sessionHandler := handlers.NewSessionHandler(sessionUsecase, log)
	tariffHandler := handlers.NewTariffHandler(tariffUsecase, log)
	walletHandler := handlers.NewWalletHandler(walletUsecase, log)
	computerHandler := handlers.NewComputerHandler(computerUsecase, log)

	// Инициализация роутера
	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware(log))

	// Регистрация маршрутов
	httpService.RegisterRoutes(r, userHandler, tariffHandler, sessionHandler, walletHandler, computerHandler)

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
		UserUsecase:     &userUsecase,
		SessionUsecase:  &sessionUsecase,
		ComputerUsecase: &computerUsecase,
		TariffUsecase:   &tariffUsecase,
		WalletUsecase:   &walletUsecase,
		Router:          r,
	}
}
