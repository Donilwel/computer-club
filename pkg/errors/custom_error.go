package errors

import "errors"

var (
	// Ошибки аутентификации и авторизации
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrForbidden          = errors.New("ошибка доступа к текущему запросу. необходима роль пользователя: админ")
	ErrWrongIDFromJWT     = errors.New("неправильный user_id в токене")
	ErrWrongRoleFromJWT   = errors.New("неправильная роль в токене")
	ErrInvalidRole        = errors.New("некорректная роль")
	ErrUserAlreadyExists  = errors.New("пользователь уже существует")
	ErrInvalidCredentials = errors.New("неверный email или пароль")
	ErrUsernameTaken      = errors.New("пользователь с таким никнеймом уже существует")
	ErrMissingToken       = errors.New("токен не найден в хедере")
	ErrWrongToken         = errors.New("токен не правильный")
	ErrTokenGeneration    = errors.New("ошибка генерации токена")

	// Ошибки при работе с пользователем
	ErrHashedPassword   = errors.New("проблема с хешированием пароля")
	ErrPasswordTooShort = errors.New("пароль должен содержать минимум 6 символов")
	ErrCreatedUser      = errors.New("ошибка при создании в базу данных пользователя")
	ErrFindUser         = errors.New("ошибка при поиске пользователя в базе данных")
	ErrNameEmpty        = errors.New("имя не может быть пустым")
	ErrPasswordEmpty    = errors.New("пароль не может быть пустым")
	ErrEmailEmpty       = errors.New("почта не может быть пустая")
	ErrRegistration     = errors.New("ошибка при регистрации пользователя")

	// Ошибки сессий
	ErrSessionNotFound              = errors.New("сессия не найдена")
	ErrSessionActive                = errors.New("у пользователя уже есть активная сессия")
	ErrInvalidSessionID             = errors.New("некорректный идентификатор сессии")
	ErrCreatedSession               = errors.New("ошибка при создании в базу данных сессии")
	ErrUpdateSession                = errors.New("ошибка обновления сессии в базу данных")
	ErrSessionAlreadyFinished       = errors.New("ошибка статус сессии уже законченный")
	ErrStatusSessionAlreadyFinished = errors.New("ошибка, запрашиваемая сессия уже закончена")
	ErrFailedStatus                 = errors.New("статус запрашиваемой сессии неверный")
	ErrCacheSession                 = errors.New("ошибка добавления в кеш")

	// Ошибки компьютеров
	ErrUpdateComputer       = errors.New("ошибка обновления компьютера в базу данных")
	ErrFindComputer         = errors.New("ошибка при поиске компьютера в базе данных")
	ErrComputerNotFound     = errors.New("компьютер не найден")
	ErrPCBusy               = errors.New("компьютер уже занят")
	ErrUpdateComputerStatus = errors.New("ошибка при обновлении статуса компьютера")

	// Ошибки тарифов
	ErrFindTariffByID  = errors.New("ошибка при поиске тарифа в базе данных по id")
	ErrInvalidTariffID = errors.New("ошибка чтения тариф id")
	ErrTariffNotFound  = errors.New("тариф не найден")
	ErrFindTariffs     = errors.New("ошибка при поиске тарифов в базе данных")
	ErrTariffsNotFound = errors.New("тарифы не найден")

	// Ошибки транзакций и кошелька
	ErrInvalidAmount       = errors.New("количество денег должно быть больше нуля")
	ErrInsufficientFunds   = errors.New("баланс меньше чем сумма которую хотят вывести")
	ErrToDeposit           = errors.New("ошибка при обновлении суммы в кошельке")
	ErrCheckBalance        = errors.New("ошибка при проверке баланса, кошелек не найден")
	ErrCheckTransaction    = errors.New("ошибка при просмотре транзакций в базе данных")
	ErrWithdraw            = errors.New("ошибка при выводе средств")
	ErrWalletAlreadyExists = errors.New("ошибка при создании кошелька: он уже существует")
	ErrCreateWallet        = errors.New("ошибка при создании кошелька в базу данных")
	ErrInvalidUserID       = errors.New("ошибка чтения id пользователя")
	ErrCreateTransaction   = errors.New("ошибка создания модели транзакции")
	ErrorTypeTransaction   = errors.New("ошибка типа транзакции")

	// Ошибки базы данных
	ErrCommitData       = errors.New("ошибка при сохранении изменений в базе данных")
	ErrStartTransaction = errors.New("ошибка при создании транзакции")

	// Ошибки кэша и Redis
	ErrDeleteRedis = errors.New("ошибка удаления данных из Редис")

	// Общие ошибки
	ErrJSONRequest = errors.New("ошибка JSON запроса")
	ErrUnexpected  = errors.New("неизвестная ошибка")
	ErrCodingaData = errors.New("ошибка обработки данных")
)
