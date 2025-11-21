# GitHub Copilot Custom Instructions

## Project Overview
null3 (pronounced as nulls, *NUHLZ*) is a web application for mood tracking, built with Go (backend) and Angular (frontend). This is a pet project focusing on simplicity and ease of development. It is not intended for production use.

## Important Notes

### Breaking Changes Policy
**Breaking changes are allowed.** This is a pet project that is not supposed to be used by others and will never reach v1. Feel free to make any breaking changes if they improve the codebase.

### Architecture Assumptions
**Never check for 32/64 bit architecture.** The application asserts 64-bit system at startup (see `backend/cmd/server/main.go` lines 18-21). Do not add additional architecture checks in the codebase.

## Technology Stack
- **Backend**: Go 1.25
  - Framework: Echo v4
  - Database: SQLite with GORM
  - Authentication: JWT tokens
- **Frontend**: Angular with Node.js 20.x

## Build and Run

### Development
```bash
# Backend (no hot-reloading)
cd backend
go run cmd/server/main.go

# Frontend
cd frontend
npm ci
ng serve
```

### Production Build
```bash
make release
PRODUCTION=true ./null3-server
```

### Testing
```bash
# Backend
cd backend
make test

# Frontend
cd frontend
npm test
```

## Configuration
Environment variables can be set in the `.env` file:
- `ADDRESS`: Backend server address (default: `localhost:8080`)
- `ENABLE_CORS`: Enable CORS (default: `false`)
- `FRONTEND_URL`: Frontend URL for CORS (default: `http://localhost:4200`)
- `PRODUCTION`: Production mode (default: `false`)
- `JWT_SECRET`: JWT signing secret (required in production)
- `JWT_EXPIRATION`: JWT expiration time (default: `24h`)
- `REFRESH_TOKEN_EXPIRATION`: Refresh token expiration (default: `168h`)
- `SECURE_COOKIES`: Enable secure cookies (default: `false`)
- `DATABASE_URL`: SQLite database URL (default: `file:null3.db?_fk=1`)
- `LOG_LEVEL`: Logging level - `debug`, `info`, `warn`, `error` (default: `info`)
- `LOG_FORMAT`: Logging format - `fancy`, `text`, `json` (default: `text`)
- `ENABLE_FRONTEND_DIST`: Serve frontend from backend (default: `false`)
- `API_URL`: API URL for frontend (default: `http://localhost:8080/api`)

## Coding Standards

### Go Backend
- Follow standard Go conventions and idioms
- Use the existing error handling patterns in the codebase
- Leverage the logging package for structured logging
- Use GORM for database operations
- Follow the module-based architecture (see `internal/` directory structure)

### Angular Frontend
- Follow Angular best practices
- Use TypeScript strictly
- Follow the component-based architecture

## Project Structure
```
null3/
├── backend/
│   ├── cmd/server/         # Application entry point
│   └── internal/           # Internal packages
│       ├── core/           # Core functionality (auth, db, server, etc.)
│       └── mood/           # Mood tracking module
└── frontend/               # Angular application
```

## Current State
- Stub user authentication is implemented (user_id: 1, login: admin, password: password)
- User management (registration, password reset) is not yet implemented
- The application is in active development

## Additional Context
- This is a learning/fun project
- Simplicity is preferred over complex enterprise patterns
- The focus is on functionality and ease of development
