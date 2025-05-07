# Система распределённого вычислителя арифметических выражений

Система состоит из двух независимых сервисов, именнованных как оркестратор и агент. Подробно структура будет описана ниже.

Стек: Go, gRPC, sqlite, Docker

Возможности:
- Регистрация и авторизация пользователя
- Прием произвольного математического выражения с операторами: +, -, *, / и скобками ()
- Просмотр всех результатов переданных выражений пользователя
- Просмотр результата выражения по его айди

Особенности:
- Кеширование результатов вычисления выражений. 
Это позволяет не пересчитывать недавно вычисленное выражение.

## Установка

1. Клонирование репозитория
```
git clone --single-branch main https://github.com/child6yo/y-lms-discalc 
cd y-lms-discalc 
```
2. Конфигурация
    - Переименуйте .env.example => .env
    - Конфигурируйе параметры
3. Запуск
```
docker-compose up --build
```

Ручной запуск крайне не рекомендуется. Однако вы можете это сделать следующим образом:

1. Запуск оркестратора
```
cd orchestrator
go run cmd/main.go
```
2. Запуск агента
```
cd agent
go run cmd/main.go
```
Работоспособность при таком запуске гарантирована не будет.

## Схема архитектуры приложения

![image](https://github.com/user-attachments/assets/768f06db-0825-45fe-b4e8-b8c16ca3f5c4)


## Примеры работы

ВНИМАНИЕ: curl тестировался в bash терминале.

### POST register
201 - пользователь создан
```
curl -X POST 'localhost:8000/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user",
  "password": "password"
}'
```

Ошибка 422 - невалидные данные:
```
curl -X POST 'localhost:8000/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user",
  "password": 123
}'
```

Ошибка 500 - внутренняя ошибка сервера:
```
curl -X POST 'localhost:8000/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user",
  "password": "password"
}'
```

### POST login
200 - успешно:
```
curl -X POST 'localhost:8000/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user",
  "password": "password"
}'
```
Возвращает JWT для дальнейшей авторизации пользователя.

400 - невалидные данные:
```
curl -X POST 'localhost:8000/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "login": "aaaa",
  "password": "password"
}'
```

500 - внутренняя ошибка сервера:
```
curl -X POST 'localhost:8000/api/v1/login' \
--header 'Content-Type: application/json' 
```

### POST calculate
201 - принято в обработку:
```
curl -X POST 'localhost:8000/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YOUR_JWT' \
--data '{
  "expression": "2+2*2"
}'
```
Возвращает айди выражения.


Ошибка 422 - невалидные данные:
```
curl -X POST 'localhost:8000/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer YOUR_JWT' \
--data '{
  "expression":
}'
```


Ошибка 500 - внутренняя ошибка сервера:
```
curl -X POST 'localhost:8000/api/v1/calculate' \
--header 'Content-Type: application/json'
--header 'Authorization: Bearer YOUR_JWT' \
```


### GET expressions/:id

200 - успешно:
```
curl -X GET 'localhost:8000/api/v1/expressions/1' \
--header 'Authorization: Bearer YOUR_JWT' \
```
У выражения может быть 3 статуса:
- Calculating... - выражение в процессе вычислений
- ERROR - произошла ошибка при вычислении
- Success - выражение успешно вычислено


404 - нет такого выражения:
```
curl -X GET 'localhost:8000/api/v1/expressions/9999' \
--header 'Authorization: Bearer YOUR_JWT' \
```

### GET expressions

200 - успешно:
```
curl -X GET 'localhost:8000/api/v1/expressions' \
--header 'Authorization: Bearer YOUR_JWT' \
```
