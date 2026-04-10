# Task M3: Rust with musl static linking for a minimal scratch image.
# Alpine Rust toolchain defaults to musl target — no extra target needed.

# ── build stage ────────────────────────────────────────────────────────
FROM rust:1.82-alpine AS builder
RUN apk add --no-cache musl-dev
WORKDIR /app

# Cache dependencies in a separate layer (rebuild only when Cargo.toml changes).
COPY Cargo.toml Cargo.lock ./
RUN mkdir src && echo "fn main() {}" > src/main.rs \
    && cargo build --release \
    && rm src/main.rs target/release/rust_service

# Build actual source (only recompiles changed files, deps stay cached).
COPY src ./src
RUN touch src/main.rs && cargo build --release

# ── runtime stage: scratch ─────────────────────────────────────────────
FROM scratch
COPY --from=builder /app/target/release/rust_service /rust_service
EXPOSE 8091
ENTRYPOINT ["/rust_service"]
