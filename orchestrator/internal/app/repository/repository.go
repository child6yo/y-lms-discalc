package repository

import (
	"database/sql"

	"github.com/child6yo/y-lms-discalc/orchestrator"
)

// Cache определяет интерфейс кэша, хранящего результаты успешно вычисленых арифметических выражений.
type Cache interface {
	// Put добавляет арифметическое выражение в кэш.
	// На вход принимает структуру арифметического выражения.
	Put(result *orchestrator.Expression)

	// Get получает значение из кэша.
	// На вход принимает ключ в виде арифметического выражения в строковом виде.
	Get(expression string) (*orchestrator.Expression, bool)
}

// Database определяет интерфейс базы данных.
type Database interface {
	// CreateUser создает пользователя в БД.
	// Принимает на вход модель пользователя,
	CreateUser(user orchestrator.User) (int, error)

	// GetUser возвращает пользователя из БД.
	// Принимает на вход логин и пароль,
	GetUser(login, password string) (*orchestrator.User, error)

	// AddExpression добавляет арифметическое выражение в БД.
	// На вход принимает айди пользователя и модель выражения.
	AddExpression(userID int, expression *orchestrator.Expression) (int, error)

	// UpdateExpression обновляет результат и статус арифметического выражения.
	// На вход принимает модель выражения.
	UpdateExpression(expression *orchestrator.Expression) error

	// GetExpressionById возвращает арифметическое выражение по его айди.
	// На вход принимает айди выражения и айди пользователя.
	GetExpressionByID(expID, userID int) (*orchestrator.Expression, error)

	// GetExpressions возвращает слайс арифметических выражений, принадлежащих пользователю.
	// На вход принимает айди пользователя.
	GetExpressions(userID int) (*[]orchestrator.Expression, error)
}

// Repository реализует репозиторий, содержащий базу данных и кеш.
type Repository struct {
	Database
	Cache
}

// mainDatabase реализует интерфейс базы данных.
// На вход принимает указатель на подключение к БД.
type mainDatabase struct {
	db *sql.DB
}

func newMainDatabase(db *sql.DB) *mainDatabase {
	return &mainDatabase{db: db}
}

// NewRepository создает новый экземпляр репозитория.
//
// Параметры:
//   - db: указатель на подключение к БД
//   - cachaCap: вместительность LRU кэша.
func NewRepository(db *sql.DB, cacheCap int) *Repository {
	return &Repository{
		Database: newMainDatabase(db),
		Cache:    newExpressionCache(cacheCap),
	}
}
