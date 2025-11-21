# Copilot Instructions for null3

## Repository Overview

**null3** (pronounced "nulls", *NUHLZ*) is a mood tracking web application built with Go backend and Angular frontend. See `README.md` for project description and features.

- **Size**: Small (~2,345 lines Go, ~1,037 lines TypeScript)
- **Languages**: Go 1.25, TypeScript (Angular 20.x), SCSS
- **Backend Framework**: Echo v4 (Go web framework)
- **Frontend Framework**: Angular 20.x with Material Design
- **Database**: SQLite with GORM ORM
- **Authentication**: JWT tokens with refresh tokens
- **Key Libraries**: Echo, GORM, golang-jwt, Angular Material, RxJS

## Project Structure

High-level directory layout (use `tree` or `find` commands to explore current structure):

```
/
├── backend/                    # Go backend application
│   ├── cmd/server/main.go     # Application entry point
│   ├── internal/
│   │   ├── core/              # Core infrastructure modules
│   │   │   ├── auth/          # JWT authentication & middleware
│   │   │   ├── db/            # Database connection & migrations
│   │   │   ├── frontend/      # Frontend static file serving
│   │   │   ├── logging/       # Structured logging
│   │   │   └── server/        # Echo server setup
│   │   └── mood/              # Mood tracking domain logic
│   ├── go.mod, go.sum         # Go module files
│   └── Makefile               # Backend build commands
├── frontend/                   # Angular frontend application
│   ├── src/
│   │   ├── app/
│   │   │   ├── core/          # Core features (auth, pages)
│   │   │   └── domains/       # Domain features (mood)
│   │   └── environments/      # Environment configs
│   ├── angular.json           # Angular configuration
│   ├── package.json           # NPM dependencies
│   ├── tsconfig.json          # TypeScript configuration
│   ├── eslint.config.js       # ESLint configuration
│   ├── proxy.conf.json        # Dev proxy to backend
│   └── Makefile               # Frontend build commands
├── Makefile                    # Root level build orchestration
├── README.md                   # Project documentation
├── .editorconfig              # Editor formatting rules
└── .gitignore                 # Git ignore patterns
```

## Build & Development Workflow

### Requirements

- **Go**: 1.25 (note: project uses Go 1.25 but 1.24+ may work)
- **Node.js**: 20.x
- **NPM**: 10.x+ (comes with Node.js)

See `README.md` for the official requirements.

### CRITICAL: Known Build Issues & Workarounds

1. **Frontend Makefile Issue**: The frontend Makefile uses `ng` directly instead of `npx ng`, which fails if Angular CLI is not globally installed.
   - **Workaround**: Always use `npx ng` when running Angular commands directly, OR install Angular CLI globally: `npm install -g @angular/cli`
   - **Affected commands**: `make lint`, `make build` in frontend directory

2. **Internet Access Required for Frontend Build**: Angular production builds fail without internet access because they try to inline Google Fonts from googleapis.com.
   - **Error**: `Inlining of fonts failed. An error has occurred while retrieving https://fonts.googleapis.com/css2?family=Roboto:wght@300;400;500&display=swap`
   - **Workaround**: Ensure internet access when building frontend in production mode, or modify src/index.html to remove external font references

3. **Existing Linting Error**: There is a known unused variable error in `frontend/src/app/core/auth/services/auth.interceptor.ts` line 61.
   - This is a pre-existing issue, not caused by new changes
   - You may fix it as part of your changes if working in that area, but it's not required

### Command Sequences That Work

#### Initial Setup (Clean Environment)
```bash
# From repository root:
make prepare              # Runs prepare for both backend and frontend
# This runs:
# - Backend: go mod tidy && go mod download
# - Frontend: npm ci
# Takes ~60 seconds on first run
```

#### Backend Development
```bash
cd backend
make prepare             # Download Go dependencies (~30 sec)
make build               # Build binary to bin/server (~30-60 sec)
make lint                # Run staticcheck and go vet (~10 sec)
make format              # Format code with go fmt
make clean               # Remove bin/ and frontend/fs/*
make clean-db            # Remove null3.db

# Run development server (no hot-reload):
go run cmd/server/main.go

# The server starts on localhost:8080 by default
```

#### Frontend Development
```bash
cd frontend
make prepare             # Run npm ci (~60 sec on clean)
# NOTE: The following commands fail with the Makefile issue
# Use npx versions instead:

npx ng lint              # Run ESLint (~10 sec)
npm run prettier         # Format code with prettier (~5 sec)
npx ng serve             # Start dev server on localhost:4200
# Dev server uses proxy.conf.json to proxy /api to localhost:8080

# For production build (requires internet):
npx ng build --configuration production
```

#### Full Project Build (From Root)
```bash
# IMPORTANT: Always run prepare before build
make prepare             # Prepare both backend and frontend (~60 sec)
make build               # Build both (backend succeeds, frontend fails with ng issue)
make lint                # Lint both (backend works, frontend fails with ng issue)
make format              # Format both
make clean               # Clean both

# For production release:
make release
# This:
# 1. Cleans backend/internal/core/frontend/fs/
# 2. Builds frontend with npx ng (NOTE: This also fails due to Makefile issue)
# 3. Copies frontend build to backend/internal/core/frontend/fs/
# 4. Builds backend binary
# 5. Copies binary to ./null3-server
```

#### Recommended Build Sequence for PRs
```bash
# 1. Prepare dependencies (always do this first)
make prepare

# 2. Build backend
cd backend && make build && cd ..

# 3. Build frontend (use npx workaround)
cd frontend && npx ng build --configuration production && cd ..
# NOTE: This requires internet access

# 4. Lint backend
cd backend && make lint && cd ..

# 5. Lint frontend (use npx workaround)
cd frontend && npx ng lint && cd ..
# Expected: 1 error in auth.interceptor.ts (pre-existing)

# 6. Format code
make format
```

### Running the Application

See `README.md` for detailed instructions on running in development and production modes.

**Key Differences from README**:
- Use `npx ng serve` instead of `ng serve` for frontend (Makefile issue workaround)
- Frontend development server proxies `/api` to `localhost:8080` via `proxy.conf.json`
- Production build requires internet access (Google Fonts inlining)

## Configuration

See `README.md` for complete list of environment variables and their descriptions.

**Key Configuration Notes**:
- Environment variables can be set in a `.env` file in repository root or backend directory
- `JWT_SECRET` required in production, auto-generated in development
- Stub user credentials (until user management implemented):
  - Username: `admin`
  - Password: `password`
  - Email: `admin@example.com`

## Code Style & Conventions

### Editor Configuration
See `.editorconfig` for complete formatting rules. Key points:
- **General**: UTF-8, LF line endings, 2-space indentation, trim trailing whitespace, final newline
- **Go files**: Tabs (width 4), max line length 100
- **Makefile**: Tabs (width 2)

### Go Code Style
- Use `go fmt` for formatting (enforced by `make format`)
- Follow standard Go conventions
- Max line length: 100 characters
- Use tabs for indentation
- Run `staticcheck` and `go vet` via `make lint`

### TypeScript/Angular Code Style
- Use Prettier for formatting: `npm run prettier`
- ESLint configuration in `eslint.config.js`
- Component selectors: `app-` prefix, kebab-case (e.g., `app-entry-card`)
- Directive selectors: `app` prefix, camelCase
- Use SCSS for styles
- TypeScript strict mode enabled

## Testing

**IMPORTANT**: This repository currently has **no test files**. There are no Go test files (*_test.go) and no Angular spec files (*.spec.ts).

- Do not attempt to run tests as none exist
- When adding new features, tests are optional but recommended
- If adding tests, follow standard Go testing conventions (package_test.go) for backend
- If adding tests for frontend, use Jasmine/Karma as configured in angular.json

## Git Workflow

- `.gitignore` excludes: node_modules, dist, .angular, bin/, *.db, .env, IDE configs, null3-server binary
- The `backend/internal/core/frontend/fs/` directory is mostly ignored (only .gitkeep is tracked)
- Do not commit build artifacts or dependencies

## Key Implementation Notes

1. **Architecture**: The backend uses a modular structure where each domain (mood) and core service (auth, db, server) is initialized via an `InitModule` function

2. **Authentication**: JWT-based with refresh tokens. The `auth` module provides middleware that extracts user context from tokens

3. **Database**: Uses GORM with SQLite. Migrations are run automatically on startup (see `internal/core/db/migrate.go`)

4. **Frontend-Backend Integration**: In production mode, the backend serves the Angular frontend from `internal/core/frontend/fs/` directory. The frontend build is copied there during `make release`

5. **API Structure**: RESTful API with endpoints under `/api/*`. Frontend development server proxies these to the backend via `proxy.conf.json`

6. **Logging**: Uses structured logging with context. The logging package provides contextual logging functions

## Validation Before Submitting PR

1. **Always** run `make prepare` first after clean checkout
2. Build backend: `cd backend && make build`
3. Lint backend: `cd backend && make lint` (should pass with no errors)
4. Build frontend: `cd frontend && npx ng build --configuration production` (requires internet)
5. Lint frontend: `cd frontend && npx ng lint` (expect 1 pre-existing error)
6. Format code: `make format`
7. **Manual Testing**: Start both backend and frontend, test login and basic mood entry operations
8. Check that no unintended files are committed (build artifacts, node_modules, etc.)

## Important: Trust These Instructions

These instructions have been validated by running all commands and observing their behavior. When working on this repository:
- **Trust the command sequences documented here** - they have been tested
- **Use the workarounds** for known issues (npx ng, internet access)
- Only search for additional information if these instructions are incomplete or you encounter errors not documented here
- The build takes time: backend build ~30-60 sec, frontend build ~60-120 sec, npm ci ~60 sec
