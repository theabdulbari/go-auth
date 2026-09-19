# go-auth - Go REST API

A scalable RESTful API built with **Go (Golang)** and **Gin**, following a clean and modular project structure. The project is designed to provide a solid foundation for building production-ready backend applications with database integration, authentication, validation, and automated testing.

## 🚀 Features

* RESTful API development with Go and Gin
* Modular project structure
* User management
* PostgreSQL database integration
* GORM ORM
* Environment-based configuration
* Request validation
* JSON API responses
* Error handling
* Middleware support
* Authentication-ready architecture
* Unit testing
* API/integration testing
* Clean and maintainable code structure

## 🛠️ Tech Stack

* **Language:** Go
* **Web Framework:** Gin
* **Database:** PostgreSQL
* **ORM:** GORM
* **Testing:** Go Testing Package
* **API:** REST
* **Version Control:** Git

## 📋 Requirements

Before running the project, make sure you have installed:

* Go 1.22+
* PostgreSQL 14+
* Git

Check your Go installation:

```bash
go version
```

Check PostgreSQL:

```bash
psql --version
```

## 📥 Installation

Clone the repository:

```bash
git clone https://github.com/your-username/your-repository.git
```

Move into the project directory:

```bash
cd your-repository
```

Install Go dependencies:

```bash
go mod download
```

Or:

```bash
go mod tidy
```

## ⚙️ Environment Configuration

Create a `.env` file in the project root:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=go_api
DB_SSLMODE=disable
```

Make sure the database exists:

```sql
CREATE DATABASE go_api;
```

## ▶️ Run the Application

Run the application using:

```bash
go run .
```

Or:

```bash
go run main.go
```

The API will start on:

```text
http://localhost:8080
```

You can test the application:

```bash
curl http://localhost:8080
```

## 🧪 Testing

Go has a built-in testing framework, so you don't need to install an additional testing package for basic tests.

### Run all tests

```bash
go test ./...
```

### Run tests with detailed output

```bash
go test -v ./...
```

### Run a specific package

```bash
go test ./controllers
```

### Run a specific test

```bash
go test -run TestCreateUser ./...
```

### Run tests with coverage

```bash
go test -cover ./...
```

### Generate a coverage report

```bash
go test ./... -coverprofile=coverage.out
```

Then view the coverage:

```bash
go tool cover -html=coverage.out
```

This will open an HTML coverage report showing which parts of the application are covered by tests.

## 🔬 Example Unit Test

Create a test file using the `_test.go` suffix.

For example:

```text
models/
├── user.go
└── user_test.go
```

Example:

```go
package models

import "testing"

func TestUserModel(t *testing.T) {
	user := User{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	if user.Name == "" {
		t.Error("user name should not be empty")
	}

	if user.Email == "" {
		t.Error("user email should not be empty")
	}
}
```

Run:

```bash
go test ./models
```

## 🌐 API Testing

You can test the API using:

* cURL
* Postman
* Insomnia
* HTTPie

Example:

```bash
curl http://localhost:8080/api/users
```

Create a user:

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password"
  }'
```

## 📁 Project Structure

```text
.
├── controllers/
│   └── user_controller.go
│
├── database/
│   └── database.go
│
├── models/
│   └── user.go
│
├── routes/
│   └── routes.go
│
├── services/
│   └── user_service.go
│
├── middleware/
│   └── auth.go
│
├── tests/
│   └── ...
│
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── main.go
```

## 🧩 Development Workflow

### 1. Install dependencies

```bash
go mod download
```

### 2. Run the application

```bash
go run .
```

### 3. Run tests

```bash
go test ./...
```

### 4. Check formatting

```bash
gofmt -w .
```

### 5. Check the project

```bash
go vet ./...
```

### 6. Run tests with race detection

```bash
go test -race ./...
```

### 7. Build the application

```bash
go build -o bin/app .
```

Run the compiled application:

```bash
./bin/app
```

## 🔍 Code Quality

Before submitting changes, run:

```bash
gofmt -w .
go vet ./...
go test ./...
```

For a more complete local check:

```bash
go test -race ./...
```

## 🔐 Security

* Never commit `.env` files containing passwords or secrets.
* Use environment variables for sensitive configuration.
* Validate all incoming API requests.
* Hash passwords before storing them.
* Use HTTPS in production.
* Apply authentication and authorization middleware where required.

Add `.env` to `.gitignore`:

```gitignore
.env
.env.*
!.env.example
```

You can provide an example configuration:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME=go_api
DB_SSLMODE=disable
```

## 🏗️ Build for Production

Build the application:

```bash
go build -o bin/app .
```

Run:

```bash
./bin/app
```

For Linux deployment:

```bash
GOOS=linux GOARCH=amd64 go build -o bin/app .
```

## 🧪 Complete Test Process

A recommended development test process is:

```bash
# Format
gofmt -w .

# Static analysis
go vet ./...

# Unit and integration tests
go test ./...

# Race-condition detection
go test -race ./...

# Test coverage
go test ./... -cover

# Build
go build -o bin/app .
```

If all commands complete successfully, the project is ready for further API testing or deployment.

## 📌 Useful Go Commands

```bash
# Initialize a project
go mod init your-module-name

# Install a dependency
go get github.com/gin-gonic/gin

# Download dependencies
go mod download

# Clean dependencies
go mod tidy

# Run application
go run .

# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Format code
gofmt -w .

# Static analysis
go vet ./...

# Build application
go build .

# View dependency information
go list -m all
```

## 🤝 Contributing

1. Fork the repository.
2. Create a feature branch:

```bash
git checkout -b feature/your-feature
```

3. Make your changes.
4. Format the code:

```bash
gofmt -w .
```

5. Run tests:

```bash
go test ./...
```

6. Commit your changes:

```bash
git add .
git commit -m "Add your feature"
```

7. Push the branch:

```bash
git push origin feature/your-feature
```

8. Open a Pull Request.

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

**Abdul Bari**

Full-Stack Developer | Senior Software Engineer | Engineering Team Lead

---

⭐ If you find this project useful, consider giving it a star.
