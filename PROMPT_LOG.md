# Prompt Log

> Лабораторная работа №11 — Docker: Python / Go / Rust  
> Студент: Никишина Евгения Александровна, группа 221131

---

## Задание М1: Dockerfile для Python-приложения (FastAPI)

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:**
«Прочитай файл "do this.txt" и выполни все задания лабораторной работы №11. Приложение FastAPI уже есть в lab10/task_m5_json_exchange/python_client. Нужен Dockerfile с многоэтапной сборкой.»

**Результат:** Claude изучил структуру lab10, скопировал app.py и test_app.py в python_service/, добавил чтение GO_SERVICE_URL из переменной окружения (для работы в Docker Compose), написал Dockerfile с двумя стадиями: builder (pip install --prefix=/deps) и runtime (python:3.12-slim).

### Итого

- Количество промптов: 1
- Что пришлось исправлять вручную: ничего
- Время: ~5 мин

---

## Задание М3: Dockerfile для Rust-приложения (Axum)

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:** (тот же промпт — всё сделано за один сеанс)

**Результат:** Создан Rust-сервис на Axum 0.7 с эндпоинтами GET /health и POST /items. Написаны 3 теста (test_health_ok, test_create_item_ok, test_create_item_invalid_json) с использованием tower::ServiceExt::oneshot и http_body_util. Dockerfile использует rust:1.82-alpine (musl) для статической линковки, итоговый образ — scratch. Реализован трюк с кешированием зависимостей (заглушка main.rs в отдельном слое).

### Итого

- Количество промптов: 1
- Что пришлось исправлять вручную: ничего
- Время: ~5 мин

---

## Задание М5: docker-compose.yml для трёх сервисов

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:** (тот же промпт — всё сделано за один сеанс)

**Результат:** docker-compose.yml поднимает go_service (:8082), python_service (:8090, передаёт GO_SERVICE_URL=http://go_service:8082, depends_on go_service) и rust_service (:8091). Все три сервиса собираются из локальных Dockerfile.

### Итого

- Количество промптов: 1
- Что пришлось исправлять вручную: ничего
- Время: ~2 мин

---

## Задание В1: Go-приложение со статической компиляцией в scratch-образе

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:** (тот же промпт — всё сделано за один сеанс)

**Результат:** go_service/Dockerfile использует golang:1.21-alpine + CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" для получения полностью статического бинарника. Финальный образ — FROM scratch, содержит только скомпилированный /go_service. Добавлен /health эндпоинт к Go-сервису и новый тест TestHealthEndpoint.

### Итого

- Количество промптов: 1
- Что пришлось исправлять вручную: ничего
- Время: ~3 мин

---

## Задание В3: CI/CD для трёх языков

### Промпт 1

**Инструмент:** Claude Code (CLI)

**Промпт:** (тот же промпт — всё сделано за один сеанс)

**Результат:** .github/workflows/docker.yml содержит 4 jobs: test-go (go test ./...), test-python (pytest), test-rust (cargo test), build-and-push (собирает и пушит образы в GHCR только при push в main, после прохождения всех тестов). Используется docker/build-push-action@v5 и docker/login-action@v3 с GITHUB_TOKEN.

### Итого

- Количество промптов: 1
- Что пришлось исправлять вручную: ничего
- Время: ~3 мин

---

## Общий итог по лабораторной работе

| Задание | Инструмент | Промптов | Время |
| --- | --- | --- | --- |
| М1 — Dockerfile Python (FastAPI) | Claude Code | 1 | ~5 мин |
| М3 — Dockerfile Rust (Axum, musl/scratch) | Claude Code | 1 | ~5 мин |
| М5 — docker-compose.yml | Claude Code | 1 | ~2 мин |
| В1 — Go static compilation (scratch) | Claude Code | 1 | ~3 мин |
| В3 — CI/CD GitHub Actions | Claude Code | 1 | ~3 мин |
| **Итого** | | **1** | **~18 мин** |

**Выводы:**
- Go с `CGO_ENABLED=0` даёт полностью статический бинарник, пригодный для scratch-образа (~10 МБ с -ldflags="-s -w").
- Rust на Alpine автоматически использует musl-toolchain, результирующий бинарник также совместим со scratch.
- Python не компилируется в статический бинарник, но многоэтапная сборка (pip install --prefix) позволяет не включать pip/setuptools в финальный образ.
- docker-compose позволяет легко управлять межсервисными URL через environment variables (GO_SERVICE_URL).
- CI/CD разделён на test-jobs (параллельно) и build-and-push (после всех тестов) — это предотвращает пуш сломанных образов.
