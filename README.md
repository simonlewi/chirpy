# Chirpy - A Social Media API Platform

## Overview
Chirpy is a RESTful API service for a Twitter-like social media platform built with Go. It provides user authentication, content management, and premium features through a clean, well-documented API.

## Features
- **User Management**
  - User registration and authentication
  - JWT-based access tokens
  - Refresh token functionality
  - Premium user upgrades (Chirpy Red)

- **Content Management**
  - Create and retrieve chirps (posts)
  - Content moderation with profanity filtering
  - Sort and filter chirps by author
  - 140-character limit enforcement

- **Security**
  - Password hashing
  - Token-based authentication
  - API key validation for webhooks
  - Token revocation support

## Installation

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 14 or higher
- Git

### Setup
1. Clone the repository:
```bash
git clone https://github.com/yourusername/chirpy.git
cd chirpy
```

2. Set up environment variables in `.env`:
```plaintext
DB_URL=your_postgres_server
PLATFORM=dev
JWT_SECRET=your_jwt_secret
POLKA_KEY=your_polka_key
```

3. Initialize the database:
```bash
psql -U postgres -c "CREATE DATABASE chirpy;"
sqlc generate
```

4. Build and run:
```bash
go build
./chirpy
```

## API Documentation

### Users
- `POST /api/users` - Create new user
- `GET /api/users` - List all users
- `POST /api/login` - Authenticate user
- `POST /api/refresh` - Refresh access token
- `POST /api/revoke` - Revoke refresh token

### Chirps
- `GET /api/chirps` - List chirps (supports filtering and sorting)
- `POST /api/chirps` - Create new chirp
- `DELETE /api/chirps/{chirpID}` - Delete chirp

### Premium Features
- `POST /api/polka/webhooks` - Handle premium user upgrades

For detailed API documentation, see:
- [Users Documentation](doc_users.md)
- [Chirps Documentation](doc_chirps.md)
- [Authentication Documentation](doc_auth.md)

## Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

## License
This project is licensed under the MIT License - see the LICENSE file for details.