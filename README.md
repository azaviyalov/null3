# null3

(pronounced as nulls, *NUHLZ*)

A web application for mood tracking, built with Go (backend) and Angular (frontend).

This is a pet project, focusing on simplicity and ease of development. Written for fun. It is not intended for production use (you can use it in production, but do so at your own risk).

## Features
- Track and manage mood entries
- Angular frontend
- RESTful Go backend
- Simple local development and build workflow

## Requirements
- Go 1.25
- Node.js 20.x

## Usage
- See the Makefile for available commands.
- Backend and frontend are managed in their respective directories.

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
   ng serve
   ```
3. Open your browser and navigate to `http://localhost:4200` to access the application.

### Production Build
1. Build the binary
    ```bash
    make release
    ```
2. Run the built binary:
   ```bash
   PRODUCTION=true ./null3-server
   ```
3. Open your browser and navigate to `http://localhost:8080` to access the application.

## Configuration
Environment variables can be set in the `.env` file.
- `ADDRESS` is used to set the backend server address. Default is `localhost:8080`.
- `ENABLE_CORS` is used to enable CORS. Default is `false`.
- `FRONTEND_URL` is used to set the URL of the frontend application (needed for CORS). Default is `http://localhost:4200`. Not applicable if `ENABLE_CORS` is set to `false`.
- `PRODUCTION` is used to set the application to production mode. Default is `false`.
- `JWT_SECRET` is used to sign JWT tokens. Default is generated randomly. Required if `PRODUCTION` is set to `true`.
- `JWT_EXPIRATION` is used to set the JWT token expiration time in seconds. Default is `24h`. Must be a positive duration.
- `REFRESH_TOKEN_EXPIRATION` is used to set the refresh token expiration time. Default is `168h` (7 days). Must be a positive duration.
- `SECURE_COOKIES` is used to enable secure cookies. Default is `false`. Set to `true` in production environments to ensure cookies are only sent over HTTPS.
- `DATABASE_URL` is used to configure the SQLite database. Default is `file:null3.db?_fk=1`.
- `LOG_LEVEL` is used to set the logging level. Default is `info`. Options are `debug`, `info`, `warn` and `error`.
- `LOG_FORMAT` is used to set the logging format. Default is `text`. Options are `fancy`, `text` and `json`.
- `ENABLE_FRONTEND_DIST` is used to enable serving the frontend from the backend. Default is `false`.
- `API_URL` is used to replace the `%%API_URL%%` placeholder in the frontend build with the actual API URL. Default is `http://localhost:8080/api`. Not applicable if `ENABLE_FRONTEND_DIST` is set to `false`.

## Stub user
Until user management is implemented, stub user is used:
- user_id: `1`.
- login: `admin`.
- password: `password`.
- email: `admin@example.com`.

## TODOs
- [ ] Implement user management (registration, login, password reset)
- [ ] Remove user data set from the environment variables
- [ ] Add more home page features (e.g., mood statistics, charts)
- [ ] Improve error handling and logging

## License
MIT
