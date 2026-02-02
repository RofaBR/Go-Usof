# Go-Usof

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

Go-Usof is a **fully functional authentication service** designed as the foundation for a Q&A platform (similar to Stack Overflow). The project showcases enterprise-level Go development with clean architecture, type-safe database operations, dual storage strategy, and external service integrations.

**Current Status:** Authentication microservice complete and deployment-ready. Q&A features planned for future development.

## Features

### Implemented

- **User Authentication**
  - Registration with email validation
  - Secure login/logout with JWT tokens (access + refresh)
  - Email verification workflow with time-limited tokens
  - Token refresh with sliding expiration
  - Password hashing with bcrypt
  - Token revocation for logout
  - OAuth2 login with Google
  - Automatic account linking for existing users

- **User Profile Management**
  - Avatar upload via Cloudinary CDN with face detection
  - Profile updates

- **Infrastructure**
  - Clean architecture (4-layer: Domain → Repository → Service → Handler)
  - Type-safe database operations with BobGen ORM
  - Dual storage: PostgreSQL (persistent) + Redis (sessions/tokens)
  - Structured JSON logging with slog
  - Graceful shutdown
  - Docker containerization
  - Database migrations

### Planned

- Questions and answers system
- Comment functionality
- Voting and reputation system
- Tag management
- Full-text search
- Unit & integration tests
- API documentation
- Authentication middleware
- CORS middleware
- Rate limiting

## Tech Stack

- **Go** 1.24 - Latest stable release
- **Gin** v1.11.0 - Web framework
- **PostgreSQL** 18 - Primary database (pgx/v5 driver)
- **Redis** - Session/token storage
- **BobGen** v0.42.0 - Type-safe ORM
- **JWT** (golang-jwt v5.3.0) - Authentication
- **bcrypt** - Password hashing
- **Cloudinary** v2.14.0 - Image CDN
- **Gomail** v2 - SMTP email
- **OAuth2** (golang.org/x/oauth2) - Google authentication
- **Docker** & **Docker Compose** - Containerization

## Quick Start

### Docker (Recommended)

1. **Clone and configure**
   ```bash
   git clone https://github.com/RofaBR/Go-Usof.git
   cd Go-Usof
   cp .env.docker .env
   # Edit .env with your credentials (Cloudinary, SMTP, JWT secrets)
   ```

2. **Start services**
   ```bash
   docker-compose up -d
   ```

3. **Verify**
   ```bash
   curl http://localhost:8080/ping
   ```

The API will be available at `http://localhost:8080`

### Local Development

1. **Setup environment**
   ```bash
   cp .env.example .env
   # Configure DATABASE_URL, REDIS_*, JWT_*, CLOUDINARY_URL, SMTP_*
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Start databases** (using Docker)
   ```bash
   docker-compose up -d postgres redis
   ```

4. **Run migrations** (optional - runs automatically on startup)
   ```bash
   psql $DATABASE_URL < db/migrations/1_user_schema.up.sql
   ```

5. **Start the server**
   ```bash
   go run cmd/usof/main.go
   ```

## API Endpoints

### Health Check
```http
GET /ping
```

### Authentication (`/api/auth`)

**Register**
```http
POST /api/auth/register
Content-Type: application/json

{
  "login": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "full_name": "John Doe"
}
```

**Login** (returns access token + httpOnly refresh cookie)
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Email Verification**
```http
GET /api/auth/verify?token=<verification_token>
```

**Refresh Token**
```http
POST /api/auth/refresh
Cookie: refresh_token=<token>
```

**Logout**
```http
POST /api/auth/logout
Cookie: refresh_token=<token>
```

### OAuth2 Authentication (`/api/auth`)

**Google Login** (redirects to Google consent page)
```http
GET /api/auth/google
```

**Google Callback** (handles Google redirect, returns tokens)
```http
GET /api/auth/google/callback?code=<auth_code>&state=<state>

Response:
{
  "access_token": "eyJhbG...",
  "expires_in": 900
}
+ httpOnly cookie: refresh_token
```

> **Note:** For new users, an account is automatically created using Google profile data.
> For existing users (matched by email), the Google account is linked.

### User Management (`/api/user`)

**Upload Avatar**
```http
POST /api/user/upload/avatar
Authorization: Bearer <access_token>
Content-Type: multipart/form-data

avatar: <file>
```

## Architecture

Go-Usof follows **Clean Architecture** with strict layer separation:

```
┌─────────────────────────────────────────┐
│   Handler Layer (HTTP Controllers)      │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│   Service Layer (Business Logic)        │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│   Repository Layer (Data Access)        │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│   Domain Layer (Entities & Interfaces)  │
└─────────────────────────────────────────┘
```

### Project Structure

```
Go-Usof/
├── cmd/
│   ├── usof/              # Application entry point
│   └── cli/               # CLI utilities
├── db/migrations/         # SQL migrations
├── internal/
│   ├── app/              # Application factory (DI)
│   ├── config/           # Configuration
│   ├── domain/           # Entities & interfaces
│   ├── dto/              # Request/response DTOs
│   ├── handler/          # HTTP controllers
│   ├── models/           # Generated ORM models
│   ├── repositories/     # Data access
│   ├── services/         # Business logic
│   ├── router/           # Routes
│   └── storage/          # Database clients
├── pkg/logger/           # Logging utilities
├── docker-compose.yaml   # Multi-service deployment
├── Dockerfile            # Production build
└── bobgen.yaml          # ORM configuration
```

## Key Technical Decisions

1. **BobGen ORM** - Type-safe SQL builder with code generation for compile-time safety and better performance than GORM

2. **Interface-Based Design** - Repositories and services defined as interfaces for easy testing and dependency inversion

3. **Dual Storage Strategy**
   - PostgreSQL for persistent data (users, profiles)
   - Redis for ephemeral data (sessions, tokens) with automatic TTL

4. **JWT with Refresh Tokens** - Access tokens (15min) + refresh tokens (30 days) in httpOnly cookies

5. **Email Verification** - Required before full account access, 24-hour expiring tokens

6. **Cloudinary Integration** - CDN image storage with face detection cropping

7. **Generated Models** - Database schema as single source of truth

## Configuration

Environment variables (see `.env.example` for full template):

```bash
# Server
PORT=8080
GIN_MODE=release

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/usof

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT
JWT_ACCESS_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret
JWT_ACCESS_TTL=15          # minutes
JWT_REFRESH_TTL=7          # days

# Email
SENDER_EMAIL=noreply@example.com
SENDER_PASSWORD=your-smtp-app-password

# Cloudinary
CLOUDINARY_URL=cloudinary://key:secret@cloud-name

# OAuth2 (Google)
OAUTH2_CLIENT_ID=your-google-client-id
OAUTH2_CLIENT_SECRET=your-google-client-secret
OAUTH2_REDIRECT_URI=http://localhost:8080/api/auth/google/callback
```

## Development

### Common Commands

```bash
# Run application
go run cmd/usof/main.go

# Build binary
go build -o bin/usof cmd/usof/main.go

# Install dependencies
go mod download

# Format code
go fmt ./...

# Run tests
go test -v ./...

# Run with coverage
go test -v -cover ./...
```

### Database Migrations

Migrations run automatically on startup. For manual execution:

```bash
# Apply migration
psql $DATABASE_URL < db/migrations/1_user_schema.up.sql

# Rollback migration
psql $DATABASE_URL < db/migrations/1_user_schema.down.sql
```

### Regenerate ORM Models

After schema changes:

```bash
go run github.com/stephenafamo/bob/gen/bobgen-psql@latest
```

## Docker Deployment

### Multi-stage Build

The `Dockerfile` uses multi-stage builds:
- **Build stage**: Full Go toolchain
- **Runtime stage**: Alpine-based (~20MB final image)

### Docker Compose Stack

Includes:
- Go-Usof API service
- PostgreSQL 18 with health checks
- Redis with persistence
- Volume management

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f usof-api

# Stop services
docker-compose down
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.