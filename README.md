# Система учета внутриофисных расходов - Backend API

REST API для системы учета расходов, разработанный на Go с использованием Gin и pgx.

## Технологии

- **Go 1.23** - язык программирования
- **Gin** - веб-фреймворк
- **pgx/v5** - драйвер PostgreSQL
- **PostgreSQL** - база данных
- **JWT** - аутентификация
- **Swagger** - документация API
- **Docker** - контейнеризация

## Структура проекта

```
curswork-trpo/
├── cmd/
│   └── app/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── handlers/
│   │   ├── handlers.go          # HTTP обработчики
│   │   ├── routes.go            # Настройка маршрутов
│   │   └── swagger.go           # Swagger модели
│   ├── middleware/
│   │   └── auth.go              # Middleware (аутентификация, CORS)
│   ├── models/
│   │   └── models.go            # Модели данных
│   ├── repository/
│   │   └── repository.go        # Репозитории для работы с БД
│   └── service/
│       └── service.go           # Бизнес-логика
├── pkg/
│   └── adapters/
│       └── postgres/
│           └── client.go        # Клиент PostgreSQL с пулом соединений
├── docs/                        # Swagger документация (генерируется)
├── docker-compose.yml           # Docker Compose конфигурация
├── Dockerfile                   # Dockerfile для сборки
├── Makefile                     # Makefile с командами
├── go.mod                       # Go модули
└── README.md                    # Документация
```

## Быстрый старт

### Предварительные требования

- Go 1.23+
- Docker и Docker Compose
- Make (опционально)

### Установка и запуск

1. **Клонируйте репозиторий**
```bash
git clone <repository-url>
cd curswork-trpo
```

2. **Установите зависимости и сгенерируйте Swagger документацию**
```bash
make install
make swagger
```

3. **Запустите через Docker Compose**
```bash
make docker-up
```

4. **Проверьте статус**
```bash
curl http://localhost:8080/health
```

5. **Откройте Swagger UI**
   Перейдите на http://localhost:8080/swagger/index.html

### Локальная разработка

1. **Установите зависимости**
```bash
make install
```

2. **Запустите PostgreSQL**
```bash
docker-compose up -d postgres
```

3. **Установите переменные окружения**
```bash
export POSTGRES_HOST=localhost
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=expense_system
export POSTGRES_PORT=5432
export JWT_SECRET=your-secret-key
```

4. **Сгенерируйте Swagger документацию**
```bash
make swagger
```

5. **Запустите приложение**
```bash
make dev
```

6. **Откройте Swagger документацию**
   http://localhost:8080/swagger/index.html

## Swagger документация

API документация доступна через Swagger UI:
- **URL**: http://localhost:8080/swagger/index.html
- **Swagger JSON**: http://localhost:8080/swagger/doc.json

### Генерация документации

```bash
# Установить swag (если еще не установлен)
go install github.com/swaggo/swag/cmd/swag@latest

# Сгенерировать документацию
make swagger

# Или вручную
swag init -g cmd/app/main.go -o docs
```

## API Endpoints

Все эндпоинты подробно описаны в Swagger UI. Краткий обзор:

### Аутентификация (`/api/auth`)

- `POST /api/auth/register` - Регистрация пользователя
- `POST /api/auth/login` - Вход в систему
- `GET /api/auth/me` - Получить текущего пользователя 🔒

### Заявки на расходы (`/api/expenses`)

- `POST /api/expenses` - Создать заявку 🔒
- `GET /api/expenses` - Получить список заявок 🔒
- `GET /api/expenses/:id` - Получить заявку по ID 🔒
- `PUT /api/expenses/:id/status` - Обновить статус заявки 🔒👔
- `GET /api/expenses/statistics` - Получить статистику 🔒👔

### Бюджет (`/api/budget`)

- `GET /api/budget/current` - Получить текущий бюджет 🔒

🔒 - Требуется аутентификация  
👔 - Только для руководства

## Архитектура базы данных

### Таблицы

#### users
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

#### expense_requests
```sql
CREATE TABLE expense_requests (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    vendor VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    employee_id INTEGER NOT NULL REFERENCES users(id),
    reviewer_id INTEGER REFERENCES users(id),
    comments TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    reviewed_at TIMESTAMP
);
```

#### budgets
```sql
CREATE TABLE budgets (
    id SERIAL PRIMARY KEY,
    year INTEGER NOT NULL,
    month INTEGER NOT NULL,
    total DECIMAL(12, 2) NOT NULL,
    spent DECIMAL(12, 2) NOT NULL DEFAULT 0,
    remaining DECIMAL(12, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(year, month)
);
```

## Примеры использования

### С помощью curl

1. **Регистрация**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@company.com",
    "password": "password123",
    "firstName": "Иван",
    "lastName": "Петров",
    "role": "employee"
  }'
```

2. **Вход**
```bash
TOKEN=$(curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@company.com",
    "password": "password123"
  }' | jq -r '.token')
```

3. **Создание заявки**
```bash
curl -X POST http://localhost:8080/api/expenses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Закупка канцтоваров",
    "category": "office-supplies",
    "amount": 3500,
    "vendor": "Комус",
    "description": "Закупка бумаги, ручек и блокнотов"
  }'
```

### Через Swagger UI

1. Откройте http://localhost:8080/swagger/index.html
2. Нажмите "Authorize"
3. Введите токен в формате: `Bearer YOUR_TOKEN`
4. Используйте интерфейс для тестирования API

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| POSTGRES_HOST | Хост PostgreSQL | localhost |
| POSTGRES_USER | Пользователь БД | postgres |
| POSTGRES_PASSWORD | Пароль БД | postgres |
| POSTGRES_DB | Имя БД | expense_system |
| POSTGRES_PORT | Порт БД | 5432 |
| JWT_SECRET | Секретный ключ для JWT | your-secret-key |
| PORT | Порт приложения | 8080 |
| GIN_MODE | Режим работы Gin | debug |

## Команды Makefile

```bash
make help              # Показать справку
make install           # Установить зависимости и swag
make build             # Собрать приложение
make run               # Запустить приложение
make dev               # Запустить в режиме разработки
make test              # Запустить тесты
make clean             # Очистить артефакты сборки
make swagger           # Сгенерировать Swagger документацию
make swagger-serve     # Сгенерировать docs и подсказать URL
make docker-up         # Запустить Docker контейнеры
make docker-down       # Остановить Docker контейнеры
make docker-rebuild    # Пересобрать и перезапустить
make docker-logs       # Показать логи Docker
make mod-tidy          # Привести в порядок модули
```

## Особенности реализации

### Подключение к базе данных

Используется пул соединений pgx/v5 для эффективной работы с PostgreSQL:

```go
type Client struct {
    pool *pgxpool.Pool
}

// Методы для работы с БД
func (c *Client) Query(ctx, query, args...)
func (c *Client) QueryRow(ctx, query, args...)
func (c *Client) Exec(ctx, query, args...)
```

Конфигурация пула:
- MaxConns: 50
- MinConns: 1
- MaxConnLifetime: 1 hour
- MaxConnIdleTime: 30 minutes
- HealthCheckPeriod: 1 minute

### Архитектура

Приложение следует чистой архитектуре:
- **Handlers** - HTTP обработчики и маршрутизация
- **Service** - бизнес-логика
- **Repository** - работа с БД (SQL запросы)
- **Models** - модели данных
- **Middleware** - промежуточное ПО

### Безопасность

- Пароли хешируются с bcrypt
- JWT токены для аутентификации
- CORS настроен для безопасного взаимодействия
- Проверка ролей на уровне middleware
- Prepared statements для защиты от SQL injection

## Разработка

### Добавление нового эндпоинта

1. Создайте модель в `internal/models/models.go`
2. Добавьте методы репозитория в `internal/repository/`
3. Реализуйте бизнес-логику в `internal/service/`
4. Создайте обработчик в `internal/handlers/`
5. Зарегистрируйте роут в `internal/handlers/routes.go`
6. Добавьте Swagger аннотации
7. Сгенерируйте документацию: `make swagger`

### Swagger аннотации

Примеры аннотаций для handlers:

```go
// @Summary Create expense request
// @Description Create a new expense request
// @Tags expenses
// @Accept json
// @Produce json
// @Param request body models.CreateExpenseRequestDTO true "Expense data"
// @Success 201 {object} models.ExpenseRequest
// @Failure 400 {object} ErrorResponse
// @Router /api/expenses [post]
// @Security BearerAuth
func (h *ExpenseHandler) CreateExpenseRequest(c *gin.Context) {
    // ...
}
```

## Тестирование

```bash
# Запустить все тесты
go test -v ./...

# С покрытием
go test -v -cover ./...

# Конкретный пакет
go test -v ./internal/service
```

## Лицензия

MIT

## Контакты

Для вопросов и предложений создайте Issue в репозитории.

## Структура проекта

```
curswork-trpo/
├── cmd/
│   └── app/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── handlers/
│   │   └── handlers.go          # HTTP обработчики
│   ├── middleware/
│   │   └── auth.go              # Middleware (аутентификация, CORS)
│   ├── models/
│   │   └── models.go            # Модели данных
│   ├── repository/
│   │   └── repository.go        # Репозитории для работы с БД
│   └── service/
│       └── service.go           # Бизнес-логика
├── pkg/
│   └── adapters/
│       └── postgres/
│           └── client.go        # Клиент PostgreSQL
├── docker-compose.yml           # Docker Compose конфигурация
├── Dockerfile                   # Dockerfile для сборки
├── Makefile                     # Makefile с командами
├── go.mod                       # Go модули
└── README.md                    # Документация
```

## Быстрый старт

### Предварительные требования

- Go 1.23+
- Docker и Docker Compose
- Make (опционально)

### Установка и запуск

1. **Клонируйте репозиторий**
```bash
git clone <repository-url>
cd curswork-trpo
```

2. **Запустите через Docker Compose**
```bash
docker-compose up -d
```

Или с использованием Makefile:
```bash
make docker-up
```

3. **Проверьте статус**
```bash
curl http://localhost:8080/health
```

### Локальная разработка

1. **Установите зависимости**
```bash
go mod download
```

2. **Запустите PostgreSQL**
```bash
docker-compose up -d postgres
```

3. **Установите переменные окружения**
```bash
export POSTGRES_HOST=localhost
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=postgres
export POSTGRES_DB=expense_system
export JWT_SECRET=your-secret-key
```

4. **Запустите приложение**
```bash
go run cmd/app/main.go
```

## API Endpoints

### Аутентификация

#### Регистрация пользователя
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "firstName": "Иван",
  "lastName": "Иванов",
  "role": "employee"
}
```

#### Вход в систему
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstName": "Иван",
    "lastName": "Иванов",
    "role": "employee"
  }
}
```

#### Получить текущего пользователя
```http
GET /api/auth/me
Authorization: Bearer <token>
```

### Заявки на расходы

#### Создать заявку
```http
POST /api/expenses
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Закупка офисной мебели",
  "category": "furniture",
  "amount": 45000,
  "vendor": "IKEA",
  "description": "Необходимо приобрести 3 рабочих стола"
}
```

#### Получить список заявок
```http
GET /api/expenses?status=all
Authorization: Bearer <token>

Query Parameters:
- status: all | pending | approved | rejected
```

#### Получить заявку по ID
```http
GET /api/expenses/{id}
Authorization: Bearer <token>
```

#### Обновить статус заявки (только для руководства)
```http
PUT /api/expenses/{id}/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "approved",
  "comments": "Одобрено. Необходимо для работы команды."
}
```

#### Получить статистику (только для руководства)
```http
GET /api/expenses/statistics
Authorization: Bearer <token>

Response:
{
  "totalPending": 45000,
  "totalApproved": 28500,
  "pendingCount": 1,
  "approvedThisMonth": 3,
  "budgetUsed": 28500,
  "budgetRemaining": 71500
}
```

### Бюджет

#### Получить текущий бюджет
```http
GET /api/budget/current
Authorization: Bearer <token>

Response:
{
  "id": 1,
  "year": 2025,
  "month": 1,
  "total": 100000,
  "spent": 28500,
  "remaining": 71500
}
```

## Модели данных

### User (Пользователь)
```go
type User struct {
    ID        uint
    Email     string
    Password  string
    FirstName string
    LastName  string
    Role      UserRole // "employee" | "management"
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### ExpenseRequest (Заявка на расход)
```go
type ExpenseRequest struct {
    ID          uint
    Title       string
    Category    string
    Amount      float64
    Vendor      string
    Description string
    Status      RequestStatus // "pending" | "approved" | "rejected"
    EmployeeID  uint
    ReviewerID  *uint
    Comments    string
    CreatedAt   time.Time
    UpdatedAt   time.Time
    ReviewedAt  *time.Time
}
```

### Budget (Бюджет)
```go
type Budget struct {
    ID        uint
    Year      int
    Month     int
    Total     float64
    Spent     float64
    Remaining float64
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| POSTGRES_HOST | Хост PostgreSQL | localhost |
| POSTGRES_USER | Пользователь БД | postgres |
| POSTGRES_PASSWORD | Пароль БД | postgres |
| POSTGRES_DB | Имя БД | expense_system |
| POSTGRES_PORT | Порт БД | 5432 |
| JWT_SECRET | Секретный ключ для JWT | your-secret-key |
| PORT | Порт приложения | 8080 |
| GIN_MODE | Режим работы Gin | debug |

## Команды Makefile

```bash
make help           # Показать справку
make build          # Собрать приложение
make run            # Запустить приложение
make test           # Запустить тесты
make clean          # Очистить артефакты сборки
make docker-up      # Запустить Docker контейнеры
make docker-down    # Остановить Docker контейнеры
make docker-rebuild # Пересобрать и перезапустить контейнеры
make docker-logs    # Показать логи Docker
make dev            # Запустить в режиме разработки
```

## Примеры использования

### Полный сценарий работы

1. **Регистрация сотрудника**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@company.com",
    "password": "password123",
    "firstName": "Иван",
    "lastName": "Петров",
    "role": "employee"
  }'
```

2. **Вход в систему**
```bash
TOKEN=$(curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "employee@company.com",
    "password": "password123"
  }' | jq -r '.token')
```

3. **Создание заявки**
```bash
curl -X POST http://localhost:8080/api/expenses \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Закупка канцтоваров",
    "category": "office-supplies",
    "amount": 3500,
    "vendor": "Комус",
    "description": "Закупка бумаги, ручек и блокнотов"
  }'
```

4. **Получение списка заявок**
```bash
curl -X GET "http://localhost:8080/api/expenses?status=pending" \
  -H "Authorization: Bearer $TOKEN"
```

5. **Одобрение заявки (руководителем)**
```bash
curl -X PUT http://localhost:8080/api/expenses/1/status \
  -H "Authorization: Bearer $MANAGER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved",
    "comments": "Одобрено"
  }'
```

## Тестирование

```bash
# Запустить все тесты
go test -v ./...

# Запустить тесты с покрытием
go test -v -cover ./...

# Запустить тесты конкретного пакета
go test -v ./internal/service
```

## Безопасность

- Пароли хешируются с использованием bcrypt
- JWT токены используются для аутентификации
- CORS настроен для безопасного взаимодействия с фронтендом
- Роли пользователей контролируют доступ к эндпоинтам

## Разработка

### Архитектура

Приложение следует чистой архитектуре:
- **Handlers** - HTTP обработчики
- **Service** - бизнес-логика
- **Repository** - работа с БД
- **Models** - модели данных
- **Middleware** - промежуточное ПО

### Добавление нового эндпоинта

1. Создайте модель в `internal/models/models.go`
2. Добавьте методы репозитория в `internal/repository/`
3. Реализуйте бизнес-логику в `internal/service/`
4. Создайте обработчик в `internal/handlers/`
5. Зарегистрируйте роут в `cmd/app/main.go`

## Лицензия

MIT

## Контакты

Для вопросов и предложений создайте Issue в репозитории.