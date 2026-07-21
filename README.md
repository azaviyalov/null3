# null3

(pronounced as nulls, *NUHLZ*)

A small personal journal for quick mood notes and longer diary entries. The API is written in Go; the browser app uses Angular.

This is a project built just for fun. It is not intended for production use.

## Features
- Track mood records
- Write diary entries in Markdown
- Link diary entries to moods with `[[mood:<id>|label]]` or `/mood-records/<id>` links
- Ignore mood-like references inside Markdown code spans and fenced code blocks
- Follow links in either direction
- Invite-only user registration
- Admin page for creating one-time invite links
- Cookie-based sessions with hashed refresh-token storage and password resets

## Requirements
- Go 1.26.5
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
   JWT_SECRET=replace-with-a-long-random-secret \
   ADMIN_PASSWORD=replace-with-a-long-random-password \
   ENABLE_FRONTEND_DIST=true ./null3-server
   ```
3. Open `http://localhost:8080`.

## Configuration
Environment variables can be set in `.env`.

- `ADDRESS`: backend listen address. Default: `localhost:8080`.
- `ENABLE_CORS`: enable CORS. Default: `false`.
- `FRONTEND_URL`: base URL for generated invite and password-reset links, and the allowed frontend origin when CORS is enabled. Default: `http://localhost:4200`.
- `JWT_SECRET`: required JWT signing key. Use a long random value; there is no default.
- `ADMIN_PASSWORD`: required password for the configuration-only administrator. Use a long random value; there is no default.
- `JWT_EXPIRATION`: JWT lifetime. Default: `24h`; must be positive.
- `REFRESH_TOKEN_EXPIRATION`: refresh-token lifetime. Default: `168h`; must be positive.
- `PASSWORD_RESET_TOKEN_EXPIRATION`: password-reset lifetime. Default: `1h`; must be positive.
- `SECURE_COOKIES`: send cookies only over HTTPS. Default: `false`.
- `DATABASE_URL`: SQLite connection string. Default: `file:null3.db?_fk=1`.
- `LOG_LEVEL`: `debug`, `info`, `warn`, or `error`. Default: `info`.
- `LOG_FORMAT`: `fancy`, `text`, or `json`. Default: `text`.
- `ENABLE_FRONTEND_DIST`: serve the embedded frontend. Default: `false`.
- `API_URL`: API URL inserted when the embedded frontend is enabled. Default: `http://localhost:8080/api`.

## Backend tests

Run unit tests without SQLite integration tests:

```bash
make test-backend-unit
```

Run the complete backend suite, including isolated SQLite integration tests:

```bash
make test-backend
```

Run the complete backend suite and print total production-code statement
coverage (the full profile is saved to `backend/coverage.out`):

```bash
make coverage-backend
```

Generate an HTML coverage report at `backend/coverage.html`:

```bash
make coverage-backend-html
```

## Administrator access

Set `ADMIN_PASSWORD=replace-with-a-long-random-password`, restart the application, and open `/admin/login`. The form accepts only the configured password. The administrator is not stored in the database.

The admin access token lasts 30 minutes and has no refresh token. After expiration, enter the password again. Changing the password requires updating the environment and restarting the application.

## Generate secrets

The optional helper below generates `JWT_SECRET` and `ADMIN_PASSWORD` and writes them to the specified env file:

```bash
cd backend
go run ./cmd/generate-secrets .env
```

The command creates the file if needed. It exits without changing the file if either variable is already present. The helper uses only the Go standard library and is not included in the release binary.

## TODOs
- [ ] Add more home page features (e.g., mood statistics, charts)
- [ ] Improve error handling and logging

## License
MIT
