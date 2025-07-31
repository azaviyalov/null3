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
- Go 1.24
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
   ./null3-server
   ```
3. Open your browser and navigate to `http://localhost:8080` to access the application.

## Configuration
Environment variables can be set in the `.env` file.
- `HOST` is used to set the backend server host. Default is `localhost:8080`.
- `ENABLE_CORS` is used to enable CORS. Default is `false`.
- `FRONTEND_URL` is used to set the URL of the frontend application (needed for CORS). Default is `http://localhost:4200`. Not applicable if `ENABLE_CORS` is set to `false`.
- `PRODUCTION` is used to set the application to production mode. Default is `false`.
- `JWT_SECRET` is used to sign JWT tokens. Default is generated randomly. Required if `PRODUCTION` is set to `true`.
- `JWT_EXPIRATION` is used to set the JWT token expiration time in seconds. Default is `24h`. Must be a positive duration.
- `SECURE_COOKIES` is used to enable secure cookies. Default is `false`. Set to `true` in production environments to ensure cookies are only sent over HTTPS.
- `DATABASE_URL` is used to configure the SQLite database. Default is `file:null3.db?_fk=1`.
- `LOG_LEVEL` is used to set the logging level. Default is `info`. Options are `debug`, `info`, `warn` and `error`.
- `LOG_FORMAT` is used to set the logging format. Default is `json`. Options are `fancy`, `text` and `json`.
- `ENABLE_FRONTEND_DIST` is used to enable serving the frontend from the backend. Default is `false`.
- `API_URL` is used to replace the `%%API_URL%%` placeholder in the frontend build with the actual API URL. Default is `http://localhost:8080/api`. Not applicable if `ENABLE_FRONTEND_DIST` is set to `false`.

These environmental variables are for development purposes only, until user management is implemented.
- The `USER_ID` variable is used to set the user ID for the user. Default is `1`.
- The `LOGIN` variable is used to set the login for the user. Default is `admin`.
- The `PASSWORD` variable is used to set the password for the user. Default is `password`.
- The `EMAIL` variable is used to set the email for the user. Default is `admin@example.com`.

## TODOs
- [ ] Implement user management (registration, login, password reset)
- [ ] Remove user data set from the environment variables
- [ ] Add JWT token refresh functionality
- [ ] Add more home page features (e.g., mood statistics, charts)
- [ ] Improve error handling and logging

## License
MIT
