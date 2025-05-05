package service

import (
	"github.com/child6yo/y-lms-discalc/orchestrator"
	"github.com/child6yo/y-lms-discalc/orchestrator/pkg/repository"
)

// Service определяет интерфейс сервиса, содержащего всю бизнес-логику.
type Service interface {
	// PostfixExpression разбивает арифметическое выражение в соответствии с обратной польской нотацией.
	// На вход принимает выражение.
	PostfixExpression(expression string) ([]string, error)

	// CulculateExpression запускает вычисление арифметического выражения
	// и возвращает его айди.
	//
	// На вход принимает айди пользователя и выражение.
	//
	// В случае, если результат такого выражения есть в кеше, сразу же создает в базе данных
	// выражение с результатом, хранящимся в кеше.
	//
	// В случае, если такого выражения в кеше нет - передает его в канал обработки,
	// запуская процесс его вычисления.
	CulculateExpression(userID int, expression string) (int, error)

	// UpdateExpression используется обработчиком выражений для обновления результата вычисления
	// в базе данных.
	//
	// В случае, если при вычислении не произошло ошибок, также записывает его в кеш.
	UpdateExpression(result *orchestrator.Expression) error

	// GetExpressioById возвращает выражение по его айди.
	// На вход принимает айди пользователя и айди выражения.
	GetExpressioByID(userID, expID int) (*orchestrator.Expression, error)

	// GetExpressions возращает слайс всех выражений пользователя.
	// На вход принимает айди пользователя.
	GetExpressions(userID int) (*[]orchestrator.Expression, error)

	// CreateUser создает нового пользователя.
	// На вход принимает модель пользователя.
	CreateUser(user orchestrator.User) (int, error)

	// GenerateToken генерирует и возвращает новый JWT.
	// На вход принимает логин и пароль пользователя.
	GenerateToken(username, password string) (string, error)

	// ParseToken обрабатывает JWT и возвращает из него айди пользователя.
	// На вход принимает JWT.
	ParseToken(accessToken string) (int, error)
}

// MainService реализует сервис.
// На вход требует указатель на репозиторий и указатель на канал обработки выражений.
type MainService struct {
	repo       *repository.Repository
	expChannel *chan *orchestrator.Expression
}

// NewService создает новый сервис.
//
// Параметры:
//   - repository: указатель на репозиторий.
//   - expChannel: указатель на канал обработки выражений.
func NewService(repository *repository.Repository, expChannel *chan *orchestrator.Expression) *MainService {
	return &MainService{repo: repository, expChannel: expChannel}
}
