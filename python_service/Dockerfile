# Task M1: multi-stage Dockerfile for the Python FastAPI service.
# Build stage installs dependencies; runtime stage contains only what's needed.

# ── build stage: install Python packages ─────────────────────────────
FROM python:3.12-slim AS builder
WORKDIR /build
COPY requirements.txt .
RUN pip install --no-cache-dir --prefix=/deps -r requirements.txt

# ── runtime stage ─────────────────────────────────────────────────────
FROM python:3.12-slim
WORKDIR /app
COPY --from=builder /deps /usr/local
COPY app.py .
RUN adduser --disabled-password --gecos "" --no-create-home appuser
USER appuser
EXPOSE 8090
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "8090"]
