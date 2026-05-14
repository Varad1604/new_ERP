import pytest
from fastapi.testclient import TestClient
import os
import sys

# Add the parent directory to sys.path to allow importing from app
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from app.main import app

client = TestClient(app)

def test_health_check():
    response = client.get("/health")
    assert response.status_code == 200
    assert response.json() == {"status": "up", "service": "ai-engine"}

def test_ai_query_stub():
    # Force provider to be stub for deterministic testing
    os.environ["AI_PROVIDER"] = "stub"

    request_data = {
        "query": "What is the total revenue for this quarter?"
    }

    response = client.post("/api/v1/ai/query", json=request_data)
    assert response.status_code == 200

    data = response.json()
    assert "Enterprise ERP Stub Response" in data["response"]
    assert "What is the total revenue" in data["response"]
    assert data["metadata"]["provider"] == "stub"
