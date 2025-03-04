# Boilme

Boilme is a powerful, feature-rich web application framework for Go, designed to streamline web development by providing a robust set of tools and libraries that work seamlessly together.

## Overview

Boilme (Latin for "swiftness") is inspired by frameworks like Laravel but built specifically for Go. It offers a comprehensive suite of functionality that modern web applications require, from database management to authentication, all within a cohesive, easy-to-use framework.

## Key Features

### 🚀 Rapid Development
- CLI tool for scaffolding applications, models, handlers, and more
- Support for multiple template engines (Go templates and Jet templates)
- Built-in authentication system
- Automatic CSRF protection

### 🔄 Database Support
- Connect to PostgreSQL, MySQL, or MariaDB
- Database migration system
- Model generation
- Query builder integration

### 💾 Caching
- Support for Redis and Badger caching
- Simple, consistent API for different cache backends

### 🔐 Security
- Built-in CSRF protection
- Session management
- Remember-me functionality
- Secure password hashing and validation

### 📧 Mailer
- HTML and plain text email templates
- Multiple mail delivery providers (SMTP, API)
- Attachments support
- Mail queue functionality

### 📁 File Storage
- Multiple filesystem support (local, S3, MinIO, SFTP, WebDAV)
- Consistent API across different storage providers

### 📝 Logging
- Structured logging
- Error tracking

### 🛠️ Utilities
- URL signing
- Form validation
- Response helpers (JSON, XML, file downloads)
- Random string generation
- File utilities

## System Requirements

- Go 1.17 or higher
- Database (PostgreSQL, MySQL, or MariaDB if needed)
- Redis (optional, for caching and sessions)

## Getting Started

### Installation

To install Boilme, you need to have Go installed on your system.

```bash
# Clone the repo (or go get it)
git clone https://github.com/user/boilme.git

# Build the CLI tool
cd boilme
make build_cli
```

### Creating a New Project

Create a new Boilme project using the CLI:

```bash
boilme new myapp
cd myapp
```

This creates a new application with the default structure.

### Project Structure

A typical Boilme project has the following structure:

```
myapp/
├── cmd/
│   └── cli/         # CLI commands
├── data/            # Database models
├── handlers/        # HTTP handlers
├── migrations/      # Database migrations
├── middleware/      # HTTP middleware
├── public/          # Static files
├── tmp/             # Temporary files and uploads
├── views/           # Templates
├── Makefile         # Build and run commands
└── .env             # Environment variables
```

### Configuration

Boilme uses environment variables for configuration. A `.env` file is created in your project root with default values. Modify these as needed:

```env
# Application
APP_NAME=myapp
APP_URL=http://localhost:4000
DEBUG=true
PORT=4000
RPC_PORT=12345

# Database
DATABASE_TYPE=postgres  # postgres, mysql, or mariadb
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASS=password
DATABASE_NAME=boilme
DATABASE_SSL_MODE=disable

# Redis (optional)
REDIS_HOST=localhost:6379
REDIS_PASSWORD=
REDIS_PREFIX=myapp

# Sessions
SESSION_TYPE=cookie  # cookie, redis, mysql, or postgres
COOKIE_NAME=myapp
COOKIE_LIFETIME=1440
COOKIE_PERSIST=true
COOKIE_SECURE=false
COOKIE_DOMAIN=localhost

# Mail
SMTP_HOST=localhost
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_ENCRYPTION=none  # none, ssl, or tls
FROM_NAME=Info
FROM_ADDRESS=info@example.com

# Templates
RENDERER=jet  # jet or go
```

### Creating Models

Generate a new model using the CLI:

```bash
boilme make model User
```

This creates a new model file in the `data` directory. Example:

```go
package data

import (
    "time"
    up "github.com/upper/db/v4"
)

// User struct
type User struct {
    ID        int       `db:"id,omitempty"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
    // Add your fields here
}

// Table returns the table name
func (t *User) Table() string {
    return "users"
}

// GetAll gets all records from the database, using upper
func (t *User) GetAll(condition up.Cond) ([]*User, error) {
    // Implementation
}

// Get gets one record from the database, by id, using upper
func (t *User) Get(id int) (*User, error) {
    // Implementation
}

// Update updates a record in the database, using upper
func (t *User) Update(m User) error {
    // Implementation
}

// Delete deletes a record from the database by id, using upper
func (t *User) Delete(id int) error {
    // Implementation
}

// Insert inserts a model into the database, using upper
func (t *User) Insert(m User) (int, error) {
    // Implementation
}
```

### Creating Migrations

Generate a new migration using the CLI:

```bash
boilme make migration create_users_table
```

This creates two migration files in the `migrations` directory: one for up (creating the table) and one for down (dropping the table).

Run migrations:

```bash
boilme migrate
```

### Creating Handlers

Generate a new handler using the CLI:

```bash
boilme make handler Home
```

This creates a new handler file in the `handlers` directory:

```go
package handlers

import (
    "net/http"
)

// Home comment goes here
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

### Authentication

Boilme provides built-in authentication functionality. Set it up using:

```bash
boilme make auth
```

This creates the necessary database tables, models, and handlers for a complete authentication system, including:

- User registration
- Login/logout
- Password reset
- Remember me functionality
- OAuth providers (GitHub, Google)

### Middleware

Boilme comes with several built-in middleware components:

- CSRF protection
- Session management
- Authentication check
- API authentication
- Remember token handling
- Maintenance mode

Register middleware in your routes.go file:

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
)

func (a *Application) routes() http.Handler {
    mux := chi.NewRouter()
    
    // Add middleware
    mux.Use(a.Middleware.CheckRemember)
    mux.Use(a.Middleware.NoSurf)
    mux.Use(a.SessionLoad)
    
    // Routes
    mux.Get("/", a.Handlers.Home)
    
    return mux
}
```

### Templates

Boilme supports two template engines: Go templates and Jet templates. Configure the template engine in your `.env` file:

```env
RENDERER=jet  # or go
```

Render templates in your handlers:

```go
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
    err := h.App.Render.Page(w, r, "home", nil, nil)
    if err != nil {
        h.App.ErrorLog.Println(err)
    }
}
```

### File Storage

Boilme provides a unified API for different file storage providers:

```go
// Upload a file to S3
err := app.UploadFile(r, "uploads", "file", app.S3)
if err != nil {
    // Handle error
}

// List files from MinIO
files, err := app.Minio.List("uploads")
if err != nil {
    // Handle error
}

// Download a file from SFTP
err := app.SFTP.Get("downloads", "file.txt")
if err != nil {
    // Handle error
}
```

### Caching

Boilme supports Redis and Badger for caching:

```go
// Set a cache value
err := app.Cache.Set("key", "value", 3600)  // With expiration in seconds
if err != nil {
    // Handle error
}

// Get a cache value
val, err := app.Cache.Get("key")
if err != nil {
    // Handle error
}

// Delete a cache value
err := app.Cache.Forget("key")
if err != nil {
    // Handle error
}

// Empty cache by pattern
err := app.Cache.EmptyByMatch("prefix*")
if err != nil {
    // Handle error
}
```

### Mailer

Boilme includes a powerful mailing system:

```go
// Send an email
msg := mailer.Message{
    To:       "recipient@example.com",
    Subject:  "Test Email",
    Template: "welcome",  // Looks for welcome.html.tmpl and welcome.plain.tmpl
    Data:     map[string]interface{}{"name": "John Doe"},
}

app.Mail.Jobs <- msg
res := <-app.Mail.Results
if res.Error != nil {
    // Handle error
}
```

### Validation

Validate form data:

```go
// Get form data
form := url.Values{}
form.Set("email", "john@example.com")
form.Set("password", "pass123")

// Create validator
validator := app.Validator(form)

// Add validation rules
validator.Required("email", "password")
validator.IsEmail("email", form.Get("email"))

// Check if valid
if !validator.Valid() {
    // Handle validation errors
    errors := validator.Errors
}
```

### API Responses

Boilme makes it easy to return JSON or XML responses:

```go
// Return JSON
err := app.WriteJSON(w, http.StatusOK, map[string]string{
    "message": "success",
})
if err != nil {
    // Handle error
}

// Return XML
err := app.WriteXML(w, http.StatusOK, myStruct)
if err != nil {
    // Handle error
}
```

### Session Management

Boilme provides session management:

```go
// Set a session value
app.Session.Put(r.Context(), "key", "value")

// Get a session value
val := app.Session.Get(r.Context(), "key")

// Remove a session value
app.Session.Remove(r.Context(), "key")

// Destroy a session
app.Session.Destroy(r.Context())

// Renew a session token
app.Session.RenewToken(r.Context())
```

### CLI Commands

Boilme comes with several built-in CLI commands:

```bash
# Show help
boilme help

# Create a new application
boilme new myapp

# Create a new model
boilme make model User

# Create a new handler
boilme make handler Home

# Create a new migration
boilme make migration create_users_table

# Run migrations
boilme migrate

# Rollback a migration
boilme migrate down

# Reset migrations
boilme migrate reset

# Create auth tables and handlers
boilme make auth

# Create session table
boilme make session

# Create a new mail template
boilme make mail welcome

# Generate a random encryption key
boilme make key

# Put the server in maintenance mode
boilme down

# Take the server out of maintenance mode
boilme up
```

## Running the Application

Run your Boilme application:

```bash
go run cmd/web/*.go
# Or using the Makefile
make run
```

## Testing

Boilme includes a testing framework:

```bash
# Run all tests
make test

# Show test coverage
make coverage

# Open test coverage in browser
make cover
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
