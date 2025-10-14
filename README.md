# Fuzzy Builder — система управления дефектами строительных объектов

Комплексное монолитное веб‑приложение для полного цикла работы с дефектами на строительных объектах: регистрация → назначение → исправление → проверка → закрытие. Репозиторий включает бэкенд (Go), фронтенд (React + TS), миграции PostgreSQL, docker‑compose и CI.

---

## 1. Цели и соответствие ТЗ

- Централизация управления дефектами на строительных объектах.
- Прозрачность процессов и статусов для инженерных и управленческих ролей.
- Ролевая модель доступа: engineer, manager, observer.
- Полный жизненный цикл дефектов: создание → назначение → работа → проверка → закрытие.
- Отчётность и аналитика: агрегаты по статусам/проектам, экспорт CSV.
- НФТ: контейнеризация, быстрое развертывание, базовая безопасность и логирование.

---

## 2. Технологический стек

- Backend: Go 1.24, Gorilla Mux (REST), pgx (PostgreSQL), JWT (HS256), bcrypt.
- DB: PostgreSQL (pgx pool), миграции *.sql.
- Frontend: React 18 + TypeScript, Vite, Ant Design, Zustand, Axios, Chart.js.
- Тесты: Go testing (unit), Vitest (+ happy‑dom) для фронта.
- Контейнеризация: Docker, docker‑compose.
- CI: GitHub Actions (build/test + docker build/push).

---

## 3. Архитектура и структура

Монолит с логическим разделением на слои и модули.

```
.
├─ cmd/
│  ├─ api/                 # точка входа API (main)
│  └─ migrator/            # утилита применения SQL‑миграций (main)
├─ internal/
│  ├─ transport/rest       # HTTP‑обработчики, роутер, middleware, swagger-заглушка
│  ├─ services             # бизнес‑логика (auth, defect)
│  ├─ repositories         # слой доступа к БД (users, projects, defects, attachments)
│  ├─ models               # доменные сущности
│  └─ database/postgresql  # pgx pool
├─ migrations/             # SQL‑миграции БД
├─ frontend/               # React + TS (AntD, Zustand, Axios, Chart.js)
├─ Dockerfile              # API Dockerfile
├─ docker-compose.yml      # compose: db, migrator, api, web
├─ run-tests.sh            # единый скрипт тестов (Go + Vitest)
└─ .github/workflows/ci.yml
```

---

## 4. Роли и доступ

- engineer — создание/обновление дефектов.
- manager — управление проектами (создание/редактирование/удаление), расширенные операции.
- observer — просмотр аналитики/отчётов.

Права проверяются на бэкенде (JWT middleware + role). На фронтенде manager‑only действия скрываются.

---

## 5. Модель данных

Основные таблицы:
- `users (id, email, password_hash, role, full_name, created_at)`
- `projects (id, name, description, created_by, created_at)`
- `defects (id, title, description, project_id, assigned_to, status, priority, due_date, created_by, created_at)`
- `attachments (id, defect_id, filename, filepath, uploaded_by, created_at)`

Статусы дефекта: `new → in_progress → on_review → closed` (правила переходов контролируются в сервисе). Приоритеты: `low | medium | high`.

### 5.1. ERD (Mermaid)

```mermaid
erDiagram
  USERS {
    bigint id PK
    text email
    text password_hash
    text role
    text full_name
    timestamptz created_at
  }

  PROJECTS {
    bigint id PK
    text name
    text description
    bigint created_by FK
    timestamptz created_at
  }

  DEFECTS {
    bigint id PK
    text title
    text description
    bigint project_id FK
    bigint assigned_to FK
    text status
    text priority
    timestamptz due_date
    bigint created_by FK
    timestamptz created_at
  }

  ATTACHMENTS {
    bigint id PK
    bigint defect_id FK
    text filename
    text filepath
    bigint uploaded_by FK
    timestamptz created_at
  }

  USERS ||--o{ PROJECTS : "created_by"
  USERS ||--o{ DEFECTS : "assigned_to"
  USERS ||--o{ DEFECTS : "created_by"
  USERS ||--o{ ATTACHMENTS : "uploaded_by"
  PROJECTS ||--o{ DEFECTS : "project_id"
  DEFECTS ||--o{ ATTACHMENTS : "defect_id"
```

Примечания:
- `users.role ∈ {engineer, manager, observer}`
- `defects.status ∈ {new, in_progress, on_review, closed}`
- `defects.priority ∈ {low, medium, high}`
- FK с пометкой `o{` означают, что связь опциональна со стороны потомка (например, `assigned_to` может быть NULL).

---

## 6. REST API (основные endpoints)

Базовый URL: `http://localhost:8080`

Аутентификация:
- `POST /api/auth/login` — вход, ответ `{ token }`.
- `POST /api/auth/register` — регистрация (в проде ограничивать).
- `GET  /api/users/me` — текущий пользователь (по Bearer).

Проекты (Bearer; manager‑only для модификаций):
- `GET    /api/projects` — список проектов.
- `POST   /api/projects` — создать проект.
- `PUT    /api/projects/{id}` — изменить проект.
- `DELETE /api/projects/{id}` — удалить проект.

Дефекты (Bearer):
- `GET    /api/defects` — список дефектов (фильтры: `status`, `priority`).
- `GET    /api/defects/{id}` — получить дефект.
- `POST   /api/defects` — создать дефект.
- `PUT    /api/defects/{id}/status` — сменить статус (валидируются переходы).
- `POST   /api/defects/{id}/attachments` — загрузить вложение (multipart).
- `GET    /api/defects/{id}/attachments` — список вложений.
- `GET    /api/attachments/{id}/download` — скачать вложение.

Отчётность:
- `GET /api/reports/defects` — экспорт CSV.
- `GET /api/reports/analytics` — агрегаты (по статусам и проектам).

Служебные:
- `GET /health` — проверка живости.
- `GET /swagger`, `/swagger/openapi.json` — заглушка OpenAPI.

---

## 7. Безопасность

- JWT (HS256), извлечение `userID`, `role` в middleware.
- bcrypt для паролей.
- RBAC: manager‑only в UI и на бэкенде.
- CORS, security headers, ограничение размера тела запроса, простой rate limit.
- Prepared statements (pgx) против SQL‑инъекций.
- Базовая валидация входных данных.

---

## 8. Frontend (React + TypeScript)

Страницы и роуты:
- `/login` — вход и загрузка `me`.
- `/dashboard` — метрики и графики (Chart.js) + фильтры по проекту/датам (локальные; можно вынести на бэкенд).
- `/projects` — таблица с поиском/сортировкой/пагинацией (клиент), создание/редактирование/удаление (manager).
- `/defects` — список, фильтры `status`/`priority`, поиск, действия (смена статуса, загрузка вложений).
- `/defects/new` — создание дефекта (выбор проекта из справочника).
- `/defects/:id` — детали дефекта (описание, смена статуса, вложения: список/загрузка/скачивание).
- `/reports` — экспорт CSV.

Инфраструктура:
- Zustand (`useAuthStore`), `ProtectedRoute`, компонент `<ManagerOnly>`.
- Ant Design, русская локаль, уведомления об ошибках (Axios interceptor).

---

## 9. Конфигурация (env)

Backend:
- `HTTP_ADDR` — адрес сервера (по умолчанию `:8080`).
- `DATABASE_URL` — строка подключения к БД (пример: `postgresql://postgres:postgres@db:5432/fuzzy?sslmode=disable`).
- `JWT_SECRET` — секрет подписи JWT (обязательно переопределить в проде).
- `UPLOAD_DIR` — каталог для вложений (по умолчанию `uploads`).

Frontend (Vite):
- `VITE_API_BASE` — базовый URL API (по умолчанию `/`; в compose Nginx проксирует на `api:8080`).

---

## 10. Быстрый старт (Docker Compose)

```bash
# сборка и запуск всех сервисов: db, migrator, api, web
docker compose build && docker compose up -d

# адреса
# API:      http://localhost:8080
# Frontend: http://localhost:5173
```

Сервисы compose:
- `db` — PostgreSQL 16 (volume `db_data`).
- `migrator` — применяет SQL‑миграции при старте.
- `api` — Golang API, volume `uploads` для файлов.
- `web` — Nginx со статикой фронта, прокси `/api`, `/swagger`, `/health` → API.

---

## 11. Локальная разработка

Требования: Go, Node.js, PostgreSQL.

```bash
# Бэкенд
export DATABASE_URL="postgresql://postgres:postgres@localhost:5432/fuzzy?sslmode=disable"
export JWT_SECRET="dev-secret-change-me"

go build -o bin/api cmd/api/main.go
./bin/api

# Миграции
go build -o bin/migrator cmd/migrator/main.go
DATABASE_URL="$DATABASE_URL" ./bin/migrator

# Фронтенд
cd frontend
npm i
npm run dev
# открыть http://localhost:5173
```

---

## 12. Тестирование

Единый скрипт запуска unit‑тестов бэкенда и фронтенда:

```bash
./run-tests.sh
```

- Backend: `internal/services/*_test.go` (bcrypt/JWT, правила переходов статусов).
- Frontend: Vitest (+ happy‑dom), примеры: `frontend/src/stores/auth.test.ts`, `frontend/src/components/ProtectedRoute.test.tsx`.

---

## 13. CI/CD

GitHub Actions (`.github/workflows/ci.yml`):
- `build-test`: go build/test + vitest.
- `docker`: сборка и публикация образов в GHCR: `api`, `migrator`, `frontend`.

---

## 14. Производительность и индексы

Рекомендуемые индексы:
- `defects(project_id)`, `defects(assigned_to)`, `defects(status)`, `defects(priority)`.

С ростом данных:
- перевести поиск/фильтры/пагинацию на серверные (SQL),
- добавить параметры к аналитике (projectId, dateFrom/dateTo),
- реализовать экспорт Excel (xlsx).

---

## 15. Дорожная карта (следующие шаги)

- Backend: серверные фильтры/пагинация (projects/defects), параметры аналитики (projectId/dateFrom/dateTo), экспорт Excel, аудит‑лог смен статуса.
- Frontend: серверная пагинация/поиск/фильтры, справочники пользователей/проектов в формах, улучшение UX (спиннеры, пустые состояния, 404/403).
- Безопасность: конфигурируемые CORS/RateLimit/Headers, унифицированный формат ошибок.
- DevOps: prod‑профиль compose (секреты, ресурсы), автобэкапы БД, мониторинг/метрики.
- Docs: полный OpenAPI/Swagger со схемами и гайдами по ролям.

---

## 16. Траблшутинг

- `DATABASE_URL` — используйте формат `postgresql://user:pass@host:port/db?sslmode=disable` (не `jdbc:`).
- 401 на фронтенде — проверьте `JWT_SECRET` и наличие токена; выполните вход через `/login`.
- Загрузка вложений — проверьте права на `UPLOAD_DIR` и наличие volume `uploads`.
- Nginx (web) не видит API — убедитесь, что контейнеры `api` и `web` запущены; проверьте `docker compose logs`.

---

© Fuzzy Builder. Проект готов к развёртыванию и дальнейшему расширению в соответствии с ТЗ.
