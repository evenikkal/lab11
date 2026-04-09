//! Rust Axum HTTP service — lab11 task M3.
//! Exposes GET /health and POST /items.

use axum::{
    routing::{get, post},
    Json, Router,
};
use serde::{Deserialize, Serialize};
use tokio::net::TcpListener;

#[derive(Serialize)]
struct HealthResponse {
    status: String,
    service: String,
}

#[derive(Deserialize, Serialize)]
struct Item {
    name: String,
    value: f64,
}

async fn health() -> Json<HealthResponse> {
    Json(HealthResponse {
        status: "ok".to_string(),
        service: "rust-service".to_string(),
    })
}

async fn create_item(Json(item): Json<Item>) -> Json<Item> {
    Json(item)
}

pub fn make_app() -> Router {
    Router::new()
        .route("/health", get(health))
        .route("/items", post(create_item))
}

#[tokio::main]
async fn main() {
    let listener = TcpListener::bind("0.0.0.0:8091").await.unwrap();
    println!("Rust service listening on 0.0.0.0:8091");
    axum::serve(listener, make_app()).await.unwrap();
}

#[cfg(test)]
mod tests {
    use super::*;
    use axum::{
        body::Body,
        http::{Request, StatusCode},
    };
    use http_body_util::BodyExt;
    use serde_json::json;
    use tower::ServiceExt;

    #[tokio::test]
    async fn test_health_ok() {
        let app = make_app();
        let req = Request::builder()
            .uri("/health")
            .body(Body::empty())
            .unwrap();
        let resp = app.oneshot(req).await.unwrap();
        assert_eq!(resp.status(), StatusCode::OK);

        let body = resp.into_body().collect().await.unwrap().to_bytes();
        let json: serde_json::Value = serde_json::from_slice(&body).unwrap();
        assert_eq!(json["status"], "ok");
        assert_eq!(json["service"], "rust-service");
    }

    #[tokio::test]
    async fn test_create_item_ok() {
        let app = make_app();
        let payload = json!({"name": "widget", "value": 9.99});
        let req = Request::builder()
            .method("POST")
            .uri("/items")
            .header("content-type", "application/json")
            .body(Body::from(payload.to_string()))
            .unwrap();
        let resp = app.oneshot(req).await.unwrap();
        assert_eq!(resp.status(), StatusCode::OK);

        let body = resp.into_body().collect().await.unwrap().to_bytes();
        let json: serde_json::Value = serde_json::from_slice(&body).unwrap();
        assert_eq!(json["name"], "widget");
        assert!((json["value"].as_f64().unwrap() - 9.99).abs() < 0.001);
    }

    #[tokio::test]
    async fn test_create_item_invalid_json() {
        let app = make_app();
        let req = Request::builder()
            .method("POST")
            .uri("/items")
            .header("content-type", "application/json")
            .body(Body::from("not json"))
            .unwrap();
        let resp = app.oneshot(req).await.unwrap();
        // axum returns 400 for unparseable JSON, 422 for schema mismatch
        assert_eq!(resp.status(), StatusCode::BAD_REQUEST);
    }
}
