# Лабораторная работа №11 — Контейнеризация мультиязычных приложений

**Студент:** Никишина Евгения Александровна  
**Группа:** 221131  
**Предмет:** Методы и технологии программирования

---

## Описание

Три HTTP-сервиса, упакованные в минимальные Docker-образы:

| Сервис | Язык / фреймворк | Порт | Образ |
|---|---|---|---|
| `go_service` | Go + Gin | 8082 | `scratch` (статическая компиляция) |
| `python_service` | Python + FastAPI | 8090 | `python:3.12-slim` (multi-stage) |
| `rust_service` | Rust + Axum | 8091 | `scratch` (musl статическая линковка) |

---

## Структура

```
lab11/
├── go_service/          # Go Gin Orders API
│   ├── main.go
│   ├── main_test.go
│   ├── go.mod / go.sum
│   └── Dockerfile       # CGO_ENABLED=0 → scratch (задание В1)
├── python_service/      # FastAPI gateway → Go service
│   ├── app.py
│   ├── test_app.py
│   ├── requirements.txt
│   └── Dockerfile       # multi-stage build (задание М1)
├── rust_service/        # Rust Axum HTTP service
│   ├── src/main.rs
│   ├── Cargo.toml
│   └── Dockerfile       # musl → scratch (задание М3)
├── docker-compose.yml   # задание М5
├── .github/workflows/
│   └── docker.yml       # CI/CD (задание В3)
└── PROMPT_LOG.md
```

---

## Запуск

### Все сервисы через Docker Compose

```bash
docker compose up --build
```

Сервисы будут доступны:
- Go: http://localhost:8082/health
- Python: http://localhost:8090/health
- Rust: http://localhost:8091/health

### Отдельный сервис

```bash
docker build -t go-service ./go_service && docker run -p 8082:8082 go-service
docker build -t python-service ./python_service && docker run -p 8090:8090 -e GO_SERVICE_URL=http://host.docker.internal:8082 python-service
docker build -t rust-service ./rust_service && docker run -p 8091:8091 rust-service
```

---

## Тесты

```bash
# Go
cd go_service && go test ./... -v

# Python
cd python_service && pip install -r requirements.txt && pytest -v

# Rust
cd rust_service && cargo test --verbose
```

---

## API

### Go service — Orders API (порт 8082)

| Метод | Путь | Описание |
|---|---|---|
| GET | `/health` | Статус сервиса |
| POST | `/orders` | Создать заказ |
| GET | `/orders/:id` | Получить заказ по ID |
| GET | `/orders` | Список всех заказов |

Пример запроса:
```bash
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d '{"customer_id":1,"items":[{"product_id":1,"product_name":"Laptop","quantity":1,"unit_price":1200}],"ship_to":{"street":"Lenina 5","city":"Moscow","country":"Russia","zip":"101000"}}'
```

### Python service — FastAPI gateway (порт 8090)

Те же эндпоинты `/orders` — валидирует через Pydantic и проксирует в Go-сервис.  
Документация: http://localhost:8090/docs

### Rust service — Items API (порт 8091)

| Метод | Путь | Описание |
|---|---|---|
| GET | `/health` | Статус сервиса |
| POST | `/items` | Создать и вернуть элемент |

```bash
curl -X POST http://localhost:8091/items \
  -H "Content-Type: application/json" \
  -d '{"name":"widget","value":9.99}'
```

---

## CI/CD

`.github/workflows/docker.yml` запускается при каждом push в `main`:

1. **test-go** — `go test ./...`
2. **test-python** — `pytest`
3. **test-rust** — `cargo test`
4. **build-and-push** — сборка и пуш образов в GHCR (только после прохождения всех тестов)

---

## Задания

| # | Задание | Статус |
|---|---|---|
| М1 | Dockerfile для Python (FastAPI), multi-stage | ✓ |
| М3 | Dockerfile для Rust (Axum), musl + scratch | ✓ |
| М5 | docker-compose.yml для трёх сервисов | ✓ |
| В1 | Go: статическая компиляция, scratch-образ | ✓ |
| В3 | CI/CD: сборка и пуш образов для трёх языков | ✓ |
