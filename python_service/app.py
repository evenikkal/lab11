"""
FastAPI gateway — lab11 task M1.

Validates incoming orders with Pydantic, forwards them to the Go Gin service,
and returns enriched responses.

Usage (local):
    uvicorn app:app --port 8090 --reload

In Docker / docker-compose the Go service URL is supplied via GO_SERVICE_URL env:
    GO_SERVICE_URL=http://go_service:8082 uvicorn app:app --host 0.0.0.0 --port 8090
"""

from __future__ import annotations

import os
from typing import List

import requests
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field

GO_SERVICE_URL = os.environ.get("GO_SERVICE_URL", "http://localhost:8082")

app = FastAPI(
    title="Orders Gateway (FastAPI)",
    description="Python gateway that validates and forwards orders to the Go Gin service",
    version="1.0.0",
)


# ── Pydantic models (mirror Go structs) ──────────────────────────────


class AddressIn(BaseModel):
    street: str = Field(..., min_length=1)
    city: str = Field(..., min_length=1)
    country: str = Field(..., min_length=1)
    zip: str = Field(..., min_length=1)


class OrderItemIn(BaseModel):
    product_id: int = Field(..., ge=1)
    product_name: str = Field(..., min_length=1)
    quantity: int = Field(..., ge=1)
    unit_price: float = Field(..., gt=0)


class OrderIn(BaseModel):
    customer_id: int = Field(..., ge=1)
    items: List[OrderItemIn] = Field(..., min_length=1)
    ship_to: AddressIn


class OrderOut(BaseModel):
    id: int
    customer_id: int
    items: List[OrderItemIn]
    ship_to: AddressIn
    total_amount: float
    status: str
    created_at: str


# ── Endpoints ────────────────────────────────────────────────────────


@app.get("/health")
def health():
    return {"status": "ok", "service": "fastapi-gateway"}


@app.post("/orders", response_model=OrderOut, status_code=201)
def create_order(order: OrderIn):
    """Validate with Pydantic, then forward to Go Gin service."""
    try:
        resp = requests.post(
            f"{GO_SERVICE_URL}/orders",
            json=order.model_dump(),
            timeout=10,
        )
    except requests.ConnectionError:
        raise HTTPException(status_code=502, detail="Go service unavailable")

    if resp.status_code != 201:
        raise HTTPException(status_code=resp.status_code, detail=resp.json())

    return resp.json()


@app.get("/orders/{order_id}", response_model=OrderOut)
def get_order(order_id: int):
    """Proxy GET to Go service."""
    try:
        resp = requests.get(
            f"{GO_SERVICE_URL}/orders/{order_id}",
            timeout=10,
        )
    except requests.ConnectionError:
        raise HTTPException(status_code=502, detail="Go service unavailable")

    if resp.status_code == 404:
        raise HTTPException(status_code=404, detail="order not found")
    if resp.status_code != 200:
        raise HTTPException(status_code=resp.status_code, detail=resp.json())

    return resp.json()


@app.get("/orders", response_model=List[OrderOut])
def list_orders():
    """Proxy GET to Go service."""
    try:
        resp = requests.get(f"{GO_SERVICE_URL}/orders", timeout=10)
    except requests.ConnectionError:
        raise HTTPException(status_code=502, detail="Go service unavailable")

    resp.raise_for_status()
    return resp.json()
