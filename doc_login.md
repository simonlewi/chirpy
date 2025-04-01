# Authentication API Documentation

## Login Endpoint
`POST /api/login`

Authenticates a user and returns access and refresh tokens.

### Request Body
```json
{
  "email": "string",
  "password": "string",
  "expires_in_seconds": number (optional)
}
```

### Response
```json
{
  "id": "uuid",
  "email": "string",
  "is_chirpy_red": boolean,
  "token": "string",
  "refresh_token": "string"
}
```

### Examples
```bash
# Basic login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "secretpass"}'

# Login with custom token expiration
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "secretpass",
    "expires_in_seconds": 3600
  }'
```

### Status Codes
- 200: Login successful
- 400: Invalid request body
- 401: Invalid credentials or incorrect password
- 405: Method not allowed (only POST accepted)
- 500: Server error (token creation failed)

### Notes
- The `token` is a JWT used for API authentication
- The `refresh_token` can be used to obtain a new access token
- `expires_in_seconds` is optional and defaults to system settings
- Successful response includes Chirpy Red membership status