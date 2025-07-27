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

## Configuration
- Environment variables can be set in the `.env` file in the `backend` directory.
- The `DATABASE_URL` variable is used to configure the SQLite database. Default is `file:null3.db?_fk=1`.
- The `PORT` variable is used to set the backend server port. Default is `8080`.
- The `ENABLE_FRONTEND_DIST` variable can be set to `true` to enable serving the frontend from the backend. Default is `false`.
- The `API_URL` variable is used in the frontend to point to the backend API. Default is `http://localhost:8080/api`.

## TODOs
- [ ] Implement user authentication

## License
MIT
