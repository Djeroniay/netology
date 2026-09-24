# Tasks API

REST API-сервис для управления списком задач (to-do), реализованный на Go с использованием стандартной библиотеки `net/http`.

Сервис работает без базы данных: задачи хранятся в памяти приложения и удаляются после остановки сервера. Все запросы и ответы используют JSON.

## Возможности

- Создание задачи
- Получение списка задач
- Получение одной задачи по идентификатору
- Полное обновление задачи
- Удаление задачи
- Проверка доступности сервиса через `/health`
- Валидация входящих JSON-данных
- Единый JSON-формат ошибок
- Потокобезопасное in-memory хранилище (`sync.RWMutex`)
- Логирование HTTP-запросов

## Стек

- Go
- Стандартная библиотека Go
- `net/http`
- `encoding/json`
- In-memory storage (`map[int]Task`)
- `sync.RWMutex`

## Структура проекта

```text
.
├── cmd/
│   └── server/
│       └── main.go               # Точка входа приложения
│
├── internal/
│   ├── handlers/
│   │   ├── tasks.go              # Обработчики задач
│   │   └── health.go             # Health-check обработчик
│   │
│   ├── models/
│   │   └── task.go               # Модель Task
│   │
│   └── storage/
│       ├── storage.go            # Интерфейс Storage
│       └── memory.go             # Реализация in-memory хранилища
│
├── go.mod
└── README.md
```

## Требования

Для запуска нужен установленный Go.

Проверить версию Go:

```powershell
go version
```

## Установка и запуск

Перейдите в корневую директорию проекта — туда, где расположены `go.mod`, папки `cmd` и `internal`.

```powershell
cd "D:\project\cursor_project\netology\Практическое задание Создание REST API‑сервиса (CRUD задач)"
```

Проверьте зависимости:

```powershell
go mod tidy
```

Запустите сервер:

```powershell
go run .\cmd\server\main.go
```

После успешного запуска в терминале появится сообщение:

```text
server listening on :8080
```

Сервер будет доступен по адресу:

```text
http://localhost:8080
```

Чтобы остановить сервер, нажмите:

```text
Ctrl + C
```

## Модель данных

Задача имеет следующий формат:

```json
{
  "id": 1,
  "title": "Изучить REST API",
  "done": false,
  "created_at": "2026-09-24T07:50:00Z"
}
```

| Поле | Тип | Описание |
|---|---|---|
| `id` | number | Уникальный идентификатор задачи. Генерируется сервером |
| `title` | string | Название задачи. Обязательное поле |
| `done` | boolean | Статус выполнения задачи |
| `created_at` | string | Время создания задачи в формате ISO8601 / RFC3339 |

При создании и обновлении задачи пользователь передаёт только `title` и `done`.

Пример тела запроса:

```json
{
  "title": "Изучить Go",
  "done": false
}
```

## Формат ошибок

Все ошибки возвращаются в JSON-формате:

```json
{
  "error": "task not found"
}
```

Примеры сообщений:

```json
{
  "error": "invalid JSON"
}
```

```json
{
  "error": "title is required"
}
```

```json
{
  "error": "invalid task ID"
}
```

```json
{
  "error": "method not allowed"
}
```

## Эндпоинты

| Метод | URL | Описание | Успешный статус |
|---|---|---|---|
| `GET` | `/health` | Проверка состояния сервиса | `200 OK` |
| `GET` | `/tasks` | Получить список задач | `200 OK` |
| `POST` | `/tasks` | Создать задачу | `201 Created` |
| `GET` | `/tasks/{id}` | Получить задачу по ID | `200 OK` |
| `PUT` | `/tasks/{id}` | Полностью обновить задачу | `200 OK` |
| `DELETE` | `/tasks/{id}` | Удалить задачу | `204 No Content` |

## Проверка сервиса

### GET /health

Проверяет, что сервер запущен и отвечает на запросы.

```powershell
curl.exe http://localhost:8080/health
```

Пример ответа — `200 OK`:

```json
{
  "status": "ok",
  "timestamp": "2026-09-24T07:50:00Z"
}
```

---

## Работа с задачами

### GET /tasks

Возвращает список всех задач.

```powershell
curl.exe http://localhost:8080/tasks
```

Пример ответа, когда задач нет — `200 OK`:

```json
[]
```

Пример ответа при наличии задач — `200 OK`:

```json
[
  {
    "id": 1,
    "title": "Изучить Go",
    "done": false,
    "created_at": "2026-09-24T07:50:00Z"
  },
  {
    "id": 2,
    "title": "Сделать практическое задание",
    "done": true,
    "created_at": "2026-09-24T08:00:00Z"
  }
]
```

---

### POST /tasks

Создаёт новую задачу.

```powershell
curl.exe -X POST http://localhost:8080/tasks `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"Изучить Go\",\"done\":false}"
```

Пример ответа — `201 Created`:

```json
{
  "id": 1,
  "title": "Изучить Go",
  "done": false,
  "created_at": "2026-09-24T07:50:00Z"
}
```

#### Ошибка: отсутствует title

```powershell
curl.exe -X POST http://localhost:8080/tasks `
  -H "Content-Type: application/json" `
  -d "{\"done\":false}"
```

Ответ — `400 Bad Request`:

```json
{
  "error": "title is required"
}
```

#### Ошибка: некорректный JSON

```powershell
curl.exe -X POST http://localhost:8080/tasks `
  -H "Content-Type: application/json" `
  -d "{title: Изучить Go}"
```

Ответ — `400 Bad Request`:

```json
{
  "error": "invalid JSON"
}
```

---

### GET /tasks/{id}

Возвращает одну задачу по её идентификатору.

```powershell
curl.exe http://localhost:8080/tasks/1
```

Пример ответа — `200 OK`:

```json
{
  "id": 1,
  "title": "Изучить Go",
  "done": false,
  "created_at": "2026-09-24T07:50:00Z"
}
```

#### Ошибка: задача не существует

```powershell
curl.exe http://localhost:8080/tasks/999
```

Ответ — `404 Not Found`:

```json
{
  "error": "task not found"
}
```

#### Ошибка: неверный идентификатор

```powershell
curl.exe http://localhost:8080/tasks/abc
```

Ответ — `400 Bad Request`:

```json
{
  "error": "invalid task ID"
}
```

---

### PUT /tasks/{id}

Полностью обновляет задачу. Поля `title` и `done` передаются в теле запроса.

```powershell
curl.exe -X PUT http://localhost:8080/tasks/1 `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"Изучить Go и REST API\",\"done\":true}"
```

Пример ответа — `200 OK`:

```json
{
  "id": 1,
  "title": "Изучить Go и REST API",
  "done": true,
  "created_at": "2026-09-24T07:50:00Z"
}
```

Поле `created_at` не меняется, поскольку отображает время первоначального создания задачи.

#### Ошибка: пустой title

```powershell
curl.exe -X PUT http://localhost:8080/tasks/1 `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"\",\"done\":true}"
```

Ответ — `400 Bad Request`:

```json
{
  "error": "title is required"
}
```

#### Ошибка: задача не найдена

```powershell
curl.exe -X PUT http://localhost:8080/tasks/999 `
  -H "Content-Type: application/json" `
  -d "{\"title\":\"Новая задача\",\"done\":false}"
```

Ответ — `404 Not Found`:

```json
{
  "error": "task not found"
}
```

---

### DELETE /tasks/{id}

Удаляет задачу по идентификатору.

```powershell
curl.exe -i -X DELETE http://localhost:8080/tasks/1
```

Пример ответа — `204 No Content`:

```text
HTTP/1.1 204 No Content
```

Тело ответа отсутствует.

#### Ошибка: задача уже удалена или не существует

```powershell
curl.exe -X DELETE http://localhost:8080/tasks/999
```

Ответ — `404 Not Found`:

```json
{
  "error": "task not found"
}
```

## Неподдерживаемый метод

Например, метод `PATCH` не реализован:

```powershell
curl.exe -X PATCH http://localhost:8080/tasks/1
```

Ответ — `405 Method Not Allowed`:

```json
{
  "error": "method not allowed"
}
```

## Логирование

Приложение выводит информацию о входящих запросах в консоль сервера.

Пример:

```text
2026/09/24 10:55:01 GET /health
2026/09/24 10:55:12 POST /tasks
2026/09/24 10:55:18 GET /tasks
2026/09/24 10:55:25 PUT /tasks/1
2026/09/24 10:55:30 DELETE /tasks/1
```

## Особенности реализации

- Данные хранятся в `map[int]Task`.
- Идентификатор задачи генерируется сервером, начиная с `1`.
- Для защиты данных при одновременных HTTP-запросах используется `sync.RWMutex`.
- Метод `GET` использует блокировку чтения `RLock`.
- Методы `POST`, `PUT` и `DELETE` используют блокировку записи `Lock`.
- После перезапуска приложения все созданные задачи удаляются, поскольку БД не используется.
- Все ответы, кроме успешного `DELETE` с кодом `204`, возвращаются с заголовком:

```text
Content-Type: application/json
```

## Статусы ответов

| Код | Значение | Когда используется |
|---|---|---|
| `200 OK` | Успешный запрос | Получение, обновление задачи, health-check |
| `201 Created` | Ресурс создан | Успешное создание задачи |
| `204 No Content` | Успешное удаление | Задача удалена |
| `400 Bad Request` | Неверный запрос | Невалидный JSON, пустой `title`, неправильный ID |
| `404 Not Found` | Ресурс не найден | Задача с указанным ID отсутствует |
| `405 Method Not Allowed` | Метод не поддерживается | Использован неразрешённый HTTP-метод |
| `500 Internal Server Error` | Внутренняя ошибка сервера | Непредвиденная ошибка приложения |