# Chirps API Documentation

## Endpoints

### Get Chirps
`GET /api/chirps`

Retrieves a list of chirps, with optional filtering and sorting.

#### Query Parameters
- `author_id` (optional): UUID of author to filter chirps by
- `sort` (optional): Sort order for chirps
  - `asc` (default): Ascending order by creation date
  - `desc`: Descending order by creation date

#### Response
```json
[
  {
    "id": "uuid",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "body": "string",
    "user_id": "uuid"
  }
]
```

#### Examples
```bash
# Get all chirps (default sorting)
curl http://localhost:8080/api/chirps

# Get chirps by author
curl "http://localhost:8080/api/chirps?author_id=a5aed2bf-7b75-4fd6-aaa3-1654ee6910bb"

# Get all chirps sorted by newest first
curl "http://localhost:8080/api/chirps?sort=desc"

# Get author's chirps sorted by newest first
curl "http://localhost:8080/api/chirps?author_id=a5aed2bf-7b75-4fd6-aaa3-1654ee6910bb&sort=desc"
```

#### Status Codes
- 200: Success
- 400: Invalid author ID
- 500: Server error