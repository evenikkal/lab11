# Prompt Log

> Лабораторная работа №11 — Docker: Python / Go / Rust  
> Студент: Никишина Евгения Александровна, группа 221131

---

## Задание М1: Dockerfile для Python-приложения (FastAPI)

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:**
«У меня есть FastAPI-приложение в lab10/task_m5_json_exchange/python_client (app.py, requirements.txt, test_app.py). Напиши для него Dockerfile с многоэтапной сборкой (multi-stage build). В первой стадии установи зависимости через pip install --prefix, во второй используй python:3.12-slim. Сервис должен запускаться на порту 8090 через uvicorn.»

**Результат:** Получила Dockerfile с builder-стадией (pip install --prefix=/deps) и runtime-стадией (python:3.12-slim, COPY --from=builder). Claude также предложил добавить чтение GO_SERVICE_URL из переменной окружения, чтобы сервис мог подключаться к Go-сервису внутри Docker-сети.

### Промпт 2

**Промпт:**
«Добавь в app.py чтение GO_SERVICE_URL из os.environ с дефолтным значением http://localhost:8082, чтобы при запуске в Docker Compose URL подставлялся автоматически.»

**Результат:** Добавлена строка `GO_SERVICE_URL = os.environ.get("GO_SERVICE_URL", "http://localhost:8082")`. Тесты по-прежнему зелёные — они мокируют requests, поэтому реальный URL не важен.

### Итого

- Количество промптов: 2
- Что пришлось исправлять вручную: ничего
- Время: ~10 мин

---

## Задание М3: Dockerfile для Rust-приложения (Axum)

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:**
«Создай простой HTTP-сервис на Rust с использованием Axum. Нужны два эндпоинта: GET /health (возвращает JSON {status: "ok", service: "rust-service"}) и POST /items (принимает {name: string, value: float} и возвращает то же самое). Добавь юнит-тесты с tower::ServiceExt::oneshot.»

**Результат:** Получила main.rs с make_app(), обработчиками health и create_item, и 3 тестами. Cargo.toml с axum 0.7, tokio, serde, serde_json, tower и http-body-util для тестов.

### Промпт 2

**Промпт:**
«Напиши Dockerfile для этого Rust-сервиса. Используй rust:1.82-alpine (musl для статической линковки) и финальный образ FROM scratch. Добавь трюк с кешированием зависимостей через заглушку main.rs в отдельном слое.»

**Результат:** Dockerfile с двумя стадиями: builder (rust:1.82-alpine, musl-dev, dep-caching trick) и scratch. Бинарник полностью статически слинкован.

### Промпт 3

**Промпт:**
«cargo test падает: test_create_item_invalid_json ожидает 422, но axum возвращает 400. Исправь тест.»

**Результат:** Исправила ожидаемый статус на 400 — axum возвращает 400 Bad Request для непарсируемого JSON и 422 для ошибок схемы. Все 3 теста зелёные.

### Итого

- Количество промптов: 3
- Что пришлось исправлять вручную: ничего (Claude сам нашёл и исправил ошибку)
- Время: ~15 мин

---

## Задание М5: docker-compose.yml для трёх сервисов

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:**
«Создай docker-compose.yml, который поднимает три сервиса: go_service (порт 8082, build: ./go_service), python_service (порт 8090, build: ./python_service, зависит от go_service, передаёт GO_SERVICE_URL=http://go_service:8082) и rust_service (порт 8091, build: ./rust_service).»

**Результат:** docker-compose.yml с тремя сервисами, depends_on и environment для межсервисного взаимодействия Python → Go.

### Итого

- Количество промптов: 1
- Что пришлось исправлять вручную: ничего
- Время: ~5 мин

---

## Задание В1: Go-приложение со статической компиляцией в scratch-образе

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:**
«Напиши Dockerfile для Go-сервиса (go_service/) с полностью статической компиляцией. Используй golang:1.21-alpine, флаги CGO_ENABLED=0 GOOS=linux и -ldflags="-s -w". Финальный образ — FROM scratch, содержащий только скомпилированный бинарник.»

**Результат:** Dockerfile с builder-стадией (go mod download, CGO_ENABLED=0 go build) и FROM scratch. Бинарник ~8 МБ против ~300 МБ у образа с Alpine.

### Промпт 2

**Промпт:**
«Добавь в Go-сервис эндпоинт GET /health, чтобы можно было проверить его работу через curl и настроить healthcheck в compose.»

**Результат:** Добавлен /health в SetupRouter() и тест TestHealthEndpoint. go test ./... — зелёный.

### Итого

- Количество промптов: 2
- Что пришлось исправлять вручную: ничего
- Время: ~10 мин

---

## Задание В3: CI/CD для трёх языков

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:**
«Создай GitHub Actions workflow (.github/workflows/docker.yml), который: 1) запускает тесты для каждого языка параллельно (go test ./..., pytest, cargo test), 2) после прохождения всех тестов собирает Docker-образы и пушит их в GHCR только при push в main. Используй docker/build-push-action@v5 и аутентификацию через GITHUB_TOKEN.»

**Результат:** Workflow с 4 jobs: test-go, test-python, test-rust (параллельно) и build-and-push (needs: все три, только на push в main). Образы тегируются как ghcr.io/{owner}/{service}:latest.

### Итого

- Количество промптов: 1
- Что пришлось исправлять вручную: ничего
- Время: ~5 мин

---

## Общий итог по лабораторной работе

| Задание | Инструмент | Промптов | Время |
| --- | --- | --- | --- |
| М1 — Dockerfile Python (FastAPI) | Claude Code | 2 | ~10 мин |
| М3 — Dockerfile Rust (Axum, musl/scratch) | Claude Code | 3 | ~15 мин |
| М5 — docker-compose.yml | Claude Code | 1 | ~5 мин |
| В1 — Go static compilation (scratch) | Claude Code | 2 | ~10 мин |
| В3 — CI/CD GitHub Actions | Claude Code | 1 | ~5 мин |
| **Итого** | | **9** | **~45 мин** |

**Выводы:**
- Go с `CGO_ENABLED=0` даёт полностью статический бинарник, пригодный для scratch-образа (~8 МБ с -ldflags="-s -w" против ~300 МБ на Alpine).
- Rust на Alpine автоматически использует musl-toolchain — результирующий бинарник совместим со scratch без дополнительных флагов.
- Python не компилируется в статический бинарник, но многоэтапная сборка (pip install --prefix) позволяет не включать pip/setuptools в финальный образ.
- Axum возвращает 400 для непарсируемого JSON и 422 для ошибок валидации схемы — важно учитывать в тестах.
- CI/CD разделён на test-jobs (параллельно) и build-and-push (после всех тестов) — предотвращает пуш сломанных образов.
