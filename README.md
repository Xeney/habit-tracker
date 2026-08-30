# Habit Tracker API

REST API для трекера привычек, написанный на Go.

## Стек

- Go 1.26+
- SQLite3
- Chi Router

## Запуск

```bash
go run cmd/main.go
```

Сервер запускается на `http://localhost:8080`.

### Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| PORT | Порт сервера | 8080 |
| DB_PATH | Путь к файлу базы данных | ./habits.db |

Пример:

```bash
PORT=3000 DB_PATH=./data.db go run cmd/main.go
```

## Структура проекта

```
habit-tracker/
├── cmd/
│   └── main.go              # Точка входа
├── internal/
│   ├── config/              # Конфигурация
│   │   └── config.go
│   ├── handlers/            # HTTP обработчики
│   │   ├── habits.go
│   │   └── health.go
│   ├── models/              # Модели данных
│   │   └── models.go
│   └── storage/             # Работа с базой данных
│       ├── db.go
│       └── habit_repo.go
├── docs/
│   └── API.md               # Документация API
├── go.mod
└── go.sum
```

## API

Подробная документация: [docs/API.md](docs/API.md)

### Быстрый старт

```bash
# Создать привычку
curl -X POST http://localhost:8080/api/habits \
  -H "Content-Type: application/json" \
  -d '{"title":"Медитация","goal_per_day":1}'

# Посмотреть все привычки
curl http://localhost:8080/api/habits

# Отметить выполнение
curl -X POST http://localhost:8080/api/habits/1/log

# Посмотреть статистику
curl http://localhost:8080/api/habits/1/stats
```

## Сборка

```bash
go build -o habit-tracker cmd/main.go
```

## Лицензия

MIT
