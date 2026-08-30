# Habit Tracker API

REST API для трекера привычек. Базовый URL: `http://localhost:8080/api`

## Аутентификация

API использует заголовок `X-User-ID` для идентификации пользователя. По умолчанию используется `user_id = 1`.

## Эндпоинты

### Health Check

```
GET /api/health
```

**Ответ:**
```json
{ "status": "ok" }
```

---

### Привычки

#### Получить все привычки

```
GET /api/habits
```

**Заголовки:**
| Заголовок | Тип | Обязательный | Описание |
|-----------|------|-------------|----------|
| X-User-ID | int | Нет | ID пользователя (по умолчанию 1) |

**Ответ 200:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "title": "Пить воду",
    "description": "2 литра в день",
    "goal_per_day": 8,
    "created_at": "2026-08-30T10:00:00Z",
    "streak": 5
  }
]
```

---

#### Создать привычку

```
POST /api/habits
```

**Тело запроса:**
```json
{
  "title": "Пить воду",
  "description": "2 литра в день",
  "goal_per_day": 8
}
```

| Поле | Тип | Обязательный | Описание |
|------|------|-------------|----------|
| title | string | Да | Название привычки |
| description | string | Нет | Описание |
| goal_per_day | int | Нет | Цель выполнений в день (по умолчанию 1) |

**Ответ 201:**
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Пить воду",
  "description": "2 литра в день",
  "goal_per_day": 8,
  "created_at": "2026-08-30T10:00:00Z",
  "streak": 0
}
```

---

#### Получить привычку по ID

```
GET /api/habits/{id}
```

**Параметры пути:**
| Параметр | Тип | Описание |
|----------|------|----------|
| id | int | ID привычки |

**Ответ 200:** Объект привычки

**Ответ 404:**
```json
{ "error": "habit not found" }
```

---

#### Обновить привычку

```
PUT /api/habits/{id}
```

**Тело запроса:**
```json
{
  "title": "Новое название",
  "description": "Новое описание",
  "goal_per_day": 5
}
```

Все поля опциональны. Отправляйте только те, которые нужно изменить.

**Ответ 200:** Обновлённый объект привычки

---

#### Удалить привычку

```
DELETE /api/habits/{id}
```

**Ответ 200:**
```json
{ "message": "habit deleted" }
```

---

#### Отметить выполнение

```
POST /api/habits/{id}/log
```

**Тело запроса (опционально):**
```json
{
  "count": 1
}
```

| Поле | Тип | Описание |
|------|------|----------|
| count | int | Количество выполнений (по умолчанию 1) |

**Ответ 200:**
```json
{ "message": "habit logged" }
```

---

#### Получить логи привычки

```
GET /api/habits/{id}/logs
```

**Ответ 200:**
```json
[
  {
    "id": 1,
    "habit_id": 1,
    "date": "2026-08-30",
    "count": 3
  }
]
```

---

#### Получить статистику

```
GET /api/habits/{id}/stats
```

**Ответ 200:**
```json
{
  "habit_id": 1,
  "total_logs": 42,
  "current_streak": 5,
  "best_streak": 14,
  "completion_rate": 85.5
}
```

| Поле | Описание |
|------|----------|
| total_logs | Общее количество записей |
| current_streak | Текущая серия дней |
| best_streak | Лучшая серия дней |
| completion_rate | Процент выполнения целей |

---

## Коды ответов

| Код | Описание |
|-----|----------|
| 200 | Успешный запрос |
| 201 | Ресурс создан |
| 400 | Неверный запрос |
| 404 | Ресурс не найден |
| 500 | Внутренняя ошибка сервера |

## Примеры запросов (curl)

```bash
curl http://localhost:8080/api/habits

curl -X POST http://localhost:8080/api/habits \
  -H "Content-Type: application/json" \
  -d '{"title":"Медитация","description":"10 минут утром","goal_per_day":1}'

curl -X POST http://localhost:8080/api/habits/1/log \
  -H "Content-Type: application/json" \
  -d '{"count":1}'

curl http://localhost:8080/api/habits/1/stats

curl -X PUT http://localhost:8080/api/habits/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"Йога","goal_per_day":2}'

curl -X DELETE http://localhost:8080/api/habits/1
```
