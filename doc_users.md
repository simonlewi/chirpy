# Users API Documentation

## Endpoints

### Create User
`POST /api/users`

Creates a new user account.

#### Request Body
```json
{
  "email": "string",
  "password": "string"
}
```

#### Response
```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "string",
  "is_chirpy_red": boolean
}
```

#### Status Codes
- 201: User created successfully
- 400: Invalid request body
- 500: Server error

### Get Users
`GET /api/users`

Retrieves a list of all users.

#### Response
```json
[
  {
    "id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "email": "string",
    "is_chirpy_red": boolean
  }
]
```

#### Example
```bash
# Create a new user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secretpass"}'

# Get all users
curl http://localhost:8080/api/users
```

#### Status Codes
- 200: Success
- 500: Server error