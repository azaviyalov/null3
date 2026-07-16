# null3

(pronounced as nulls, *NUHLZ*)

A small personal journal for quick mood notes and longer diary entries. The API is written in Go; the browser app uses Angular.

This is a project built just for fun. It is not intended for production use.

## Features
- Track mood records
- Write diary entries in Markdown
- Link diary entries to moods with `[[mood:<id>|label]]`
- Follow links in either direction
- Invite-only user registration
- Admin page for creating one-time invite links
- Cookie-based sessions and password resets

## Requirements
- Go 1.26.3
- Node.js 24.15.x
- npm 11.14.x

## Project Structure
- Backend `internal/core` contains infrastructure and shared runtime concerns such as database, logging, HTTP server setup, and frontend asset serving.
- Backend `internal/domain` contains feature logic such as `account`, `session`, `admin`, and `journal`.
- Frontend `src/app/core` contains shared app utilities and static app-level pages such as `about`.
- Frontend `src/app/domains` contains feature domains such as `account`, `session`, `admin`, `dashboard`, and `journal`.
- Journal pages use `/mood-records` and `/diary-entries`; their REST endpoints are grouped under `/api/journal/mood-records` and `/api/journal/diary-entries`.

## Running the Application

### Development
1. Start the backend server (no hot-reloading):
   ```bash
   cd backend
   go run cmd/server/main.go
   ```
2. Start the frontend development server:
   ```bash
   cd frontend
   npm ci
   npm start
   ```
3. Open `http://localhost:4200`.

### Production Build
1. Build the binary
    ```bash
    make release
    ```
2. Run the built binary:
   ```bash
   PRODUCTION=true ./null3-server
   ```
3. Open `http://localhost:8080`.

## Configuration
Environment variables can be set in `.env`.

- `ADDRESS`: backend listen address. Default: `localhost:8080`.
- `ENABLE_CORS`: enable CORS. Default: `false`.
- `FRONTEND_URL`: allowed frontend origin when CORS is enabled. Default: `http://localhost:4200`.
- `PRODUCTION`: production mode. Default: `false`.
- `JWT_SECRET`: JWT signing key. Generated at startup outside production; required in production.
- `JWT_EXPIRATION`: JWT lifetime. Default: `24h`; must be positive.
- `REFRESH_TOKEN_EXPIRATION`: refresh-token lifetime. Default: `168h`; must be positive.
- `PASSWORD_RESET_TOKEN_EXPIRATION`: password-reset lifetime. Default: `1h`; must be positive.
- `SECURE_COOKIES`: send cookies only over HTTPS. Default: `false`.
- `DATABASE_URL`: SQLite connection string. Default: `file:null3.db?_fk=1`.
- `LOG_LEVEL`: `debug`, `info`, `warn`, or `error`. Default: `info`.
- `LOG_FORMAT`: `fancy`, `text`, or `json`. Default: `text`.
- `ENABLE_FRONTEND_DIST`: serve the embedded frontend. Default: `false`.
- `API_URL`: API URL inserted when the embedded frontend is enabled. Default: `http://localhost:8080/api`.

## Seeded admin account
The application seeds a single admin account on startup if user `1` does not exist:
- user_id: `1`.
- login: `admin`.
- password: `password`.
- email: `admin@example.com`.

Use this account only on the separate admin login page at `/admin/login`. Regular user accounts are created through invite links generated from the admin area.

## TODOs
- [ ] Add more home page features (e.g., mood statistics, charts)
- [ ] Improve error handling and logging

## License
MIT
