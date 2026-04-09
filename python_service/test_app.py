"""
Tests for the FastAPI Orders Gateway (lab11 task M1).
All calls to the Go service are mocked — no running servers required.
"""

import unittest
from unittest.mock import patch, MagicMock

from fastapi.testclient import TestClient

from app import app

client = TestClient(app)

SAMPLE_GO_RESPONSE = {
    "id": 1,
    "customer_id": 101,
    "items": [
        {"product_id": 1, "product_name": "Laptop", "quantity": 1, "unit_price": 1200.00},
        {"product_id": 2, "product_name": "Mouse", "quantity": 2, "unit_price": 25.50},
    ],
    "ship_to": {
        "street": "Lenina 5",
        "city": "Moscow",
        "country": "Russia",
        "zip": "101000",
    },
    "total_amount": 1251.00,
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z",
}

VALID_ORDER = {
    "customer_id": 101,
    "items": [
        {"product_id": 1, "product_name": "Laptop", "quantity": 1, "unit_price": 1200.00},
        {"product_id": 2, "product_name": "Mouse", "quantity": 2, "unit_price": 25.50},
    ],
    "ship_to": {
        "street": "Lenina 5",
        "city": "Moscow",
        "country": "Russia",
        "zip": "101000",
    },
}


class TestHealth(unittest.TestCase):
    def test_health(self):
        resp = client.get("/health")
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(resp.json()["service"], "fastapi-gateway")


class TestCreateOrder(unittest.TestCase):

    @patch("app.requests.post")
    def test_create_order_success(self, mock_post):
        mock_resp = MagicMock()
        mock_resp.status_code = 201
        mock_resp.json.return_value = SAMPLE_GO_RESPONSE
        mock_post.return_value = mock_resp

        resp = client.post("/orders", json=VALID_ORDER)

        self.assertEqual(resp.status_code, 201)
        self.assertEqual(resp.json()["id"], 1)
        self.assertAlmostEqual(resp.json()["total_amount"], 1251.00)

    def test_create_order_pydantic_validates_missing_items(self):
        bad_order = {"customer_id": 101, "ship_to": VALID_ORDER["ship_to"]}
        resp = client.post("/orders", json=bad_order)
        self.assertEqual(resp.status_code, 422)

    def test_create_order_pydantic_validates_bad_customer_id(self):
        bad_order = {**VALID_ORDER, "customer_id": 0}
        resp = client.post("/orders", json=bad_order)
        self.assertEqual(resp.status_code, 422)

    def test_create_order_pydantic_validates_negative_price(self):
        bad_order = {
            **VALID_ORDER,
            "items": [{"product_id": 1, "product_name": "X", "quantity": 1, "unit_price": -5}],
        }
        resp = client.post("/orders", json=bad_order)
        self.assertEqual(resp.status_code, 422)

    def test_create_order_pydantic_validates_empty_address(self):
        bad_order = {
            **VALID_ORDER,
            "ship_to": {"street": "", "city": "M", "country": "R", "zip": "1"},
        }
        resp = client.post("/orders", json=bad_order)
        self.assertEqual(resp.status_code, 422)

    @patch("app.requests.post")
    def test_create_order_go_unavailable(self, mock_post):
        import requests as req
        mock_post.side_effect = req.ConnectionError("refused")

        resp = client.post("/orders", json=VALID_ORDER)
        self.assertEqual(resp.status_code, 502)


class TestGetOrder(unittest.TestCase):

    @patch("app.requests.get")
    def test_get_order_success(self, mock_get):
        mock_resp = MagicMock()
        mock_resp.status_code = 200
        mock_resp.json.return_value = SAMPLE_GO_RESPONSE
        mock_get.return_value = mock_resp

        resp = client.get("/orders/1")
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(resp.json()["customer_id"], 101)

    @patch("app.requests.get")
    def test_get_order_not_found(self, mock_get):
        mock_resp = MagicMock()
        mock_resp.status_code = 404
        mock_resp.json.return_value = {"error": "order not found"}
        mock_get.return_value = mock_resp

        resp = client.get("/orders/999")
        self.assertEqual(resp.status_code, 404)


class TestListOrders(unittest.TestCase):

    @patch("app.requests.get")
    def test_list_orders(self, mock_get):
        mock_resp = MagicMock()
        mock_resp.status_code = 200
        mock_resp.json.return_value = [SAMPLE_GO_RESPONSE]
        mock_resp.raise_for_status = MagicMock()
        mock_get.return_value = mock_resp

        resp = client.get("/orders")
        self.assertEqual(resp.status_code, 200)
        self.assertEqual(len(resp.json()), 1)


if __name__ == "__main__":
    unittest.main()
