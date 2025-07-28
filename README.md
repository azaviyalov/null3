# null3

(pronounced as nulls, *NUHLZ*)

A full-stack web application for mood tracking, built with Go (backend) and Angular (frontend).

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
- The `DATABASE_URL` variable is used to configure the SQLite database. Default is `file:null3.db?_fk=1`.
- The `PORT` variable is used to set the backend server port. Default is `8080`.
- The `LOG_LEVEL` variable is used to set the logging level. Default is `info`. Options are `debug`, `info`, `warn` and `error`.
- The `ENABLE_CORS` variable can be set to `true` to enable CORS. Default is `false`.
- The `ENABLE_FRONTEND_DIST` variable can be set to `true` to enable serving the frontend from the backend. Default is `false`.
- The `API_URL` variable is used to replace the `%%API_URL%%` placeholder in the frontend build with the actual API URL. Default is `http://localhost:8080/api`.

## TODOs
- [ ] Implement user authentication

## License
MIT
