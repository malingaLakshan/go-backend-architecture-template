# API Documentation

Base URL: `http://localhost:8080`

---

## Health Check

### `GET /health`

Returns the server status.

**Response (200 OK):**

```json
{
    "success": true,
    "message": "Server is running",
    "data": null
}
```

---

## Customers

### `POST /api/customers`

Create a new customer.

**Request Body:**

```json
{
    "name": "Jane Doe",
    "email": "jane.doe@example.com",
    "phone": "+1-555-0100"
}
```

**Response (201 Created):**

```json
{
    "success": true,
    "message": "Customer created successfully",
    "data": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "name": "Jane Doe",
        "email": "jane.doe@example.com",
        "phone": "+1-555-0100",
        "status": "ACTIVE",
        "created_at": "2025-01-15T10:30:00Z",
        "updated_at": "2025-01-15T10:30:00Z"
    }
}
```

**Error Response (400 Bad Request):**

```json
{
    "success": false,
    "message": "Validation failed",
    "error": "email is required"
}
```

**Error Response (409 Conflict):**

```json
{
    "success": false,
    "message": "Conflict",
    "error": "a customer with this email already exists"
}
```

---

### `GET /api/customers`

List all customers. Results are ordered by creation date (newest first).

**Response (200 OK):**

```json
{
    "success": true,
    "message": "Customers retrieved successfully",
    "data": [
        {
            "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
            "name": "Jane Doe",
            "email": "jane.doe@example.com",
            "phone": "+1-555-0100",
            "status": "ACTIVE",
            "created_at": "2025-01-15T10:30:00Z",
            "updated_at": "2025-01-15T10:30:00Z"
        },
        {
            "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "name": "John Smith",
            "email": "john.smith@example.com",
            "phone": "+1-555-0200",
            "status": "ACTIVE",
            "created_at": "2025-01-14T09:00:00Z",
            "updated_at": "2025-01-14T09:00:00Z"
        }
    ]
}
```

---

### `GET /api/customers/{id}`

Get a single customer by ID.

**Example:** `GET /api/customers/f47ac10b-58cc-4372-a567-0e02b2c3d479`

**Response (200 OK):**

```json
{
    "success": true,
    "message": "Customer retrieved successfully",
    "data": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "name": "Jane Doe",
        "email": "jane.doe@example.com",
        "phone": "+1-555-0100",
        "status": "ACTIVE",
        "created_at": "2025-01-15T10:30:00Z",
        "updated_at": "2025-01-15T10:30:00Z"
    }
}
```

**Error Response (404 Not Found):**

```json
{
    "success": false,
    "message": "Not found",
    "error": "customer not found"
}
```

---

### `PUT /api/customers/{id}`

Update an existing customer. All fields except `status` are required. If `status` is omitted, the existing status is preserved.

**Example:** `PUT /api/customers/f47ac10b-58cc-4372-a567-0e02b2c3d479`

**Request Body:**

```json
{
    "name": "Jane Doe-Smith",
    "email": "jane.smith@example.com",
    "phone": "+1-555-0101",
    "status": "INACTIVE"
}
```

**Response (200 OK):**

```json
{
    "success": true,
    "message": "Customer updated successfully",
    "data": {
        "id": "f47ac10b-58cc-4372-a567-0e02b2c3d479",
        "name": "Jane Doe-Smith",
        "email": "jane.smith@example.com",
        "phone": "+1-555-0101",
        "status": "INACTIVE",
        "created_at": "2025-01-15T10:30:00Z",
        "updated_at": "2025-01-15T14:00:00Z"
    }
}
```

**Error Response (400 Bad Request):**

```json
{
    "success": false,
    "message": "Validation failed",
    "error": "name is required"
}
```

**Error Response (404 Not Found):**

```json
{
    "success": false,
    "message": "Not found",
    "error": "customer not found"
}
```

---

### `DELETE /api/customers/{id}`

Delete a customer by ID.

**Example:** `DELETE /api/customers/f47ac10b-58cc-4372-a567-0e02b2c3d479`

**Response (200 OK):**

```json
{
    "success": true,
    "message": "Customer deleted successfully",
    "data": null
}
```

**Error Response (404 Not Found):**

```json
{
    "success": false,
    "message": "Not found",
    "error": "customer not found"
}
```

---

## Error Codes Summary

| HTTP Status | Meaning             | When                                    |
|-------------|----------------------|-----------------------------------------|
| 200         | OK                   | Successful read, update, or delete      |
| 201         | Created              | Successful create                       |
| 400         | Bad Request          | Invalid JSON or validation failure      |
| 404         | Not Found            | Customer ID does not exist              |
| 409         | Conflict             | Duplicate email                         |
| 500         | Internal Server Error| Unexpected server error                 |

---

## Testing with cURL

### Create a customer

```bash
curl -X POST http://localhost:8080/api/customers \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe", "email": "jane@example.com", "phone": "+1-555-0100"}'
```

### List all customers

```bash
curl http://localhost:8080/api/customers
```

### Get a customer by ID

```bash
curl http://localhost:8080/api/customers/<customer-id>
```

### Update a customer

```bash
curl -X PUT http://localhost:8080/api/customers/<customer-id> \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Smith", "email": "jane.smith@example.com", "phone": "+1-555-0101", "status": "INACTIVE"}'
```

### Delete a customer

```bash
curl -X DELETE http://localhost:8080/api/customers/<customer-id>
```
