## Система распределённого вычислителя арифметических выражений

Еще один проект из программы яндекс лицея. Написан на чистом Go, без сторонних библиотек.
Система состоит из двух независимых сервисов, именнованных как оркестратор и агент. Подробно структура будет описана ниже.

Возможности:
- Прием произвольного математического выражения с операторами: +, -, *, / и скобками ()
- Просмотр всех результатов переданных выражений
- Просмотр результата выражения по его айди

Каждое из переданных выражений вычисляется независимо. Кроме того, вычислением каждого из выражений занимается несколько потоков в составе агента.

## Установка

1. Клонирование репозитория
```
git clone https://github.com/child6yo/y-lms-discalc 
cd y-lms-discalc 
```
2. Конфигурация
   - Выполняется в файле docker-compose, подробно о конфигурации будет описано ниже
3. Запуск
```
docker-compose up
```
4. UI
Пользовательский интерфейс доступен по адресу:
```
http://localhost:8000
```

## Конфигурация

Вся конфигурация выполняется, изменяя передаваемые переменные окружения через docker-compose.
- Для оркестратора:
    - TIME_ADDITION_MS — время на выполнение операции сложения в миллисекундах
    - TIME_SUBTRACTION_MS — время на выполнение операции вычитания в миллисекундах
    - TIME_MULTIPLICATIONS_MS — время на выполнение операции умножения в миллисекундах
    - TIME_DIVISIONS_MS — время на выполнение операции деления в миллисекундах
- Для агента:
    - COMPUTING_POWER - количество потоков, одновременно вычисляющих пришедшие задачи
 
## Примеры работы

ВНИМАНИЕ: curl тестировался в bash терминале.

### POST calculate
201 - принято в обработку:
```
curl -X POST 'localhost:8000/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
Возвращает айди выражения.


Ошибка 422 - невалидные данные:
```
curl -X POST 'localhost:8000/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression":
}'
```


Ошибка 500 - внутренняя ошибка сервера:
```
curl -X POST 'localhost:8000/api/v1/calculate' \
--header 'Content-Type: application/json'
```


### GET expressions/:id

200 - успешно:
```
curl -X GET 'localhost:8000/api/v1/expressions/1'
```
У выражения может быть 3 статуса:
- Calculating... - выражение в процессе вычислений
- ERROR - произошла ошибка при вычислении
- Success - выражение успешно вычислено


404 - нет такого выражения:
```
curl -X GET 'localhost:8000/api/v1/expressions/9999'
```

### GET expressions

200 - успешно:
```
curl -X GET 'localhost:8000/api/v1/expressions'
```
Из-за особенностей внутреннего устройства и в целях оптимизации полученный список НЕ будет упорядочен по айди.

#
## Структура проекта

Как уже упоминалось, проект разделен на два сервиса:
- Окрестратор
```
orchestrator
├── cmd
|   └── main.go # запуск сервера
├── pkg
|   ├── handler
|   |   ├── calc.go # содержит хендлеры эндпоинтов для пользователя
|   |   ├── handler.go # общая логика хендлеров
|   |   └── internal.go # содержит хендлеры для внутренних эндпоинтов системы
|   ├── processor
|   |   └── processor.go # содержит код, запускающий процесс для каждего пришедшего выражения. процессор руководит делением и отправкой задач.
|   ├── service
|   |   └── service.go # содержит логику перевода в обратную польскую нотацию для отправки стека в процессор
└── models.go # содержит все использующиеся модели
```

- Агент
```
agent
├── cmd
|   └── main.go # запуск рабочих
├── pkg
|   ├── worker
|   |   └── worker.go # содержит функцию независимых "рабочих", обращающихся к серверу за задачами и обрабатывающими задачи
|   ├── service 
|   └── └── service.go # содержит логику вычислений.
└── models.go # содержит все использующиеся модели
```

# Ниже представлена УПРОЩЕННАЯ модель работы системы
![изображение](https://github.com/user-attachments/assets/608915f9-42fc-41ae-83ed-5138b623a0bd)


