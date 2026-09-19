# go-auth Go REST API

A scalable RESTful API built with **Go (Golang)** and **Gin**, following a clean and modular project structure. It provides a solid foundation for backend applications with SQLite database integration, authentication, validation, and automated testing.

## 🚀 Features

- RESTful API development with Go and Gin
- Modular project structure
- User management (CRUD) + authentication
- SQLite database integration (pure-Go driver — no CGO required)
- GORM ORM
- Environment-based configuration (`.env`)
- Request validation
- Consistent JSON API responses
- Centralized error handling
- JWT authentication middleware
- Password hashing with bcrypt
- Unit and integration tests
- Clean, maintainable code structure

## 🛠️ Tech Stack

- **Language:** Go
- **Web Framework:** Gin
- **Database:** SQLite
- **SQLite Driver:** `github.com/glebarez/sqlite` (pure Go, no CGO)
- **ORM:** GORM
- **Auth:** JWT (`github.com/golang-jwt/jwt/v5`) + bcrypt
- **Testing:** Go's built-in `testing` package
- **API Style:** REST
- **Version Control:** Git

## 📋 Requirements

Make sure the following are installed:

- Go 1.22+
- Git

> **No PostgreSQL or C compiler needed.** The SQLite driver is pure Go, so it works with `CGO_ENABLED=0`.

Verify your setup:

```bash
go version
```

## 📥 Installation

Clone the repository and move into the project directory:

```bash
git clone https://github.com/your-username/your-repository.git
cd your-repository
```

Install dependencies:

```bash
go mod download
# or
go mod tidy
```

## ⚙️ Environment Configuration

Create a `.env` file in the project root:

```env
APP_ENV=development
APP_PORT=8080

DB_PATH=users.db

JWT_SECRET=change-me-in-production
JWT_EXPIRY_HOURS=24
```

The SQLite database file (`users.db`) is created automatically on first run. No manual database setup required.

## ▶️ Running the Application

```bash
go run .
# or
go run main.go

# If you ever hit a CGO issue:
CGO_ENABLED=0 go run .
```

The API will be available at:

```text
http://localhost:8080
```

Quick check:

```bash
curl http://localhost:8080
```

## 🔐 Authentication Flow

1. **Register** a user → `POST /auth/register`
2. **Login** → `POST /auth/login` → returns a JWT token
3. **Access protected routes** by sending the token:

```
Authorization: Bearer <token>
```

## 🌐 API Endpoints

| Method | Endpoint         | Auth | Description              |
|--------|------------------|------|--------------------------|
| GET    | `/`              | No   | Health check             |
| POST   | `/auth/register` | No   | Create a new user        |
| POST   | `/auth/login`    | No   | Login and get JWT token  |
| GET    | `/api/users`     | Yes  | List all users           |
| GET    | `/api/profile`   | Yes  | Get current user profile |

## 🧪 API Testing

You can test the API with **cURL**, **Postman**, **Insomnia**, or **HTTPie**.

```bash
# Register
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"secret123"}'

# Login (capture the token)
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"secret123"}'

# Protected route
curl http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 🧪 Testing

Go ships with a built-in testing framework, so no extra tools are required for basic tests.

### Run all tests

```bash
go test ./...
```

### Verbose output

```bash
go test -v ./...
```

### Run a specific package

```bash
go test ./handlers
```

### Run a specific test

```bash
go test -run TestRegister ./...
```

### Coverage

```bash
go test -cover ./...

# Generate an HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Example Unit Test

Create a file with the `_test.go` suffix, e.g. `utils/password_test.go`:

```go
package utils

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	hashed, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("failed to hash: %v", err)
	}

	if !CheckPassword(hashed, "secret123") {
		t.Error("expected password to match")
	}

	if CheckPassword(hashed, "wrongpassword") {
		t.Error("expected password NOT to match")
	}
}
```

Run it:

```bash
go test ./utils
```

## 📁 Project Structure

```text
.
├── controllers/
│   └── user_controller.go
├── database/
│   └── db.go
├── handlers/
│   ├── auth.go
│   └── user.go
├── middleware/
│   └── auth.go
├── models/
│   └── user.go
├── routes/
│   └── routes.go
├── services/
│   └── user_service.go
├── utils/
│   ├── jwt.go
│   └── password.go
├── tests/
│   └── ...
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── main.go
```

## 🧩 Development Workflow

1. **Install dependencies**

   ```bash
   go mod download
   ```

2. **Run the application**

   ```bash
   go run .
   ```

3. **Run tests**

   ```bash
   go test ./...
   ```

4. **Format code**

   ```bash
   gofmt -w .
   ```

5. **Static analysis**

   ```bash
   go vet ./...
   ```

6. **Race detection**

   ```bash
   go test -race ./...
   ```

7. **Build**

   ```bash
   go build -o bin/app .
   ./bin/app
   ```

## 🔍 Code Quality

Before committing or opening a PR:

```bash
gofmt -w .
go vet ./...
go test -race ./...
```

## 🔐 Security

- Never commit `.env` files containing secrets.
- Use environment variables for sensitive configuration (especially `JWT_SECRET`).
- Validate all incoming requests.
- Hash passwords with bcrypt before storing them.
- Use the same "Invalid credentials" message for wrong email *and* wrong password (prevents user enumeration).
- Serve over HTTPS in production.
- Apply authentication and authorization middleware where required.

Add the following to `.gitignore`:

```gitignore
.env
.env.*
!.env.example

# SQLite
*.db
*.db-journal
*.db-wal

# Build
bin/
```

Provide an `.env.example` with placeholder values:

```env
APP_ENV=development
APP_PORT=8080

DB_PATH=users.db

JWT_SECRET=
JWT_EXPIRY_HOURS=24
```

## 🏗️ Build for Production

```bash
# Local build (no CGO needed)
CGO_ENABLED=0 go build -o bin/app .
./bin/app

# Linux build (cross-compile)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/app .
```

> **Tip:** Because the SQLite driver is pure Go, you get a fully static binary — perfect for Docker scratch images and simple deployment.

## 🐳 Optional: Docker

```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/app .
EXPOSE 8080
CMD ["./app"]
```

## 🧪 Complete Pre-Deploy Check

```bash
gofmt -w .
go vet ./...
go test ./...
go test -race ./...
go test ./... -cover
CGO_ENABLED=0 go build -o bin/app .
```

If all commands complete successfully, the project is ready for further API testing or deployment.

## 📌 Useful Go Commands

```bash
go mod init your-module-name            # Initialize a module
go get github.com/gin-gonic/gin         # Add a dependency
go get github.com/glebarez/sqlite       # Pure-Go SQLite driver
go get github.com/golang-jwt/jwt/v5     # JWT
go get golang.org/x/crypto/bcrypt       # bcrypt
go mod download                         # Download dependencies
go mod tidy                             # Clean up dependencies
go run .                                # Run the app
go test ./...                           # Run all tests
go test -v ./...                        # Verbose test output
gofmt -w .                              # Format code
go vet ./...                            # Static analysis
CGO_ENABLED=0 go build .                # Build (static binary)
go list -m all                          # List dependencies
```

## 🤝 Contributing

1. Fork the repository.
2. Create a feature branch:

   ```bash
   git checkout -b feature/your-feature
   ```

3. Make your changes.
4. Format and test:

   ```bash
   gofmt -w .
   go test ./...
   ```

5. Commit and push:

   ```bash
   git add .
   git commit -m "Add your feature"
   git push origin feature/your-feature
   ```

6. Open a Pull Request.

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

**Abdul Bari**
Full-Stack Developer | Senior Software Engineer | Engineering Team Lead

---

⭐ If you find this project useful, consider giving it a star.
