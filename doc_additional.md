# Additional API Endpoints Documentation

## Authentication & User Management

### Refresh Token
`POST /api/refresh`
Refreshes an expired access token using a valid refresh token.

#### Request Header
```
Authorization: Bearer <refresh_token>
```

#### Response
```json
{
  "token": "string"
}
```

### Revoke Token
`POST /api/revoke`
Revokes a refresh token to prevent further use.

#### Request Header
```
Authorization: Bearer <refresh_token>
```

#### Status Code
- 200: Token revoked successfully

## Content Validation

### Validate Chirp
`POST /api/validate_chirp`
Validates chirp content for profanity and length constraints.

#### Request Body
```json
{
  "body": "string"
}
```

#### Response
```json
{
  "cleaned_body": "string",
  "valid": boolean
}
```

## Admin Endpoints

### Metrics
`GET /admin/metrics`
Retrieves server metrics and statistics.

#### Response
```json
{
  "fileserver_hits": number
}
```

### Reset
`POST /admin/reset`
Resets server metrics (admin only).

## Health Check

### Health Check
`GET /api/healthz`
Checks API health status.

#### Response
- 200: API is healthy