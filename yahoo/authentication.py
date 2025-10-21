import hmac
import hashlib
import time
import json
import requests
# ========================================
# Configuration (provided by Yahoo! Auction Bridge Service)
# ========================================
API_KEY = "buyship-service-001"
SECRET_KEY = "your-secret-key-here"
BASE_URL = "https://yahoo-auction-bridge.example.com"
# ========================================
# Authentication Function
# ========================================
def make_authenticated_request(method, path, body=None):
    """
    Make an authenticated request to Yahoo! Auction Bridge Service
    Args:
        method: HTTP method ("GET", "POST", "PUT", "DELETE")
        path: API path (e.g., "/api/v1/auction/item")
        body: Request body (dict or None)
    Returns:
        Response object
    """
    # 1. Generate timestamp
    timestamp = str(int(time.time()))
    # 2. Convert body to JSON string (NO SPACES)
    body_str = json.dumps(body, separators=(',', ':')) if body else ""
    # 3. Send request
    url = BASE_URL + path
    response = requests.request(
        method=method,
        url=url,
        data=body_str if body else None
    )
    return response

# ========================================
# Usage Examples
# ========================================
# Example 1: Get product information (GET with query parameters)
response = make_authenticated_request(
    method="GET",
    path="/api/v1/auction/item?aID=abc123",
    body=None
)
print(response.json())
# Example 2: Buy-out purchase (POST with body)
response = make_authenticated_request(
    method="POST",
    path="/api/v1/auction/purchase",
    body={
        "aID": "abc123",
        "Quantity": 1
    }
)
print(response.json())

# Example 3: Product search (GET with query parameters)
response = make_authenticated_request(
    method="GET",
    path="/api/v1/auction/search?query=laptop&category=2084005502",
    body={}
)
print(response.json())