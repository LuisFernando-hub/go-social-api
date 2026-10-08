# Go Social API

A simple social media backend built with Go, Gin, GORM, and PostgreSQL. This project exposes REST endpoints for user registration, authentication, post creation, listing, and basic interaction actions such as like and dislike.

The API is organized using a modular internal package structure and uses JWT-based authentication for protected routes.

## Overview

This project was designed as a small social application backend where users can:

- create an account
- log in securely
- publish posts
- view the list of posts
- interact with posts using LIKE or DESLIKE actions

It follows a lightweight layered architecture with:

- handlers for HTTP layer
- services for business logic
- repositories for persistence
- models for database entities
- middleware for authentication

## Features

- User registration with password hashing using bcrypt
- JWT login flow
- Protected routes with middleware validation
- PostgreSQL persistence with GORM
- Auto-migration of database models at startup
- Post creation and listing
- Like/dislike interactions on posts
- Health endpoint for quick server validation
- Environment-based configuration via .env or OS environment variables

## Tech Stack

- Go
- Gin Web Framework
- GORM ORM
- PostgreSQL
- JWT (golang-jwt)
- bcrypt for password hashing
- dotenv support via godotenv

## Project Structure

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── handler.go
│   │   └── service.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── postgres.go
│   ├── middleware/
│   │   └── auth/
│   │       └── middleware.go
│   ├── models/
│   │   ├── post.go
│   │   ├── post_interaction.go
│   │   └── user.go
│   ├── post/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── post_interaction/
│   │   ├── handler.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── types/
│   │   └── interaction.go
│   └── user/
│       ├── handler.go
│       ├── repository.go
│       └── service.go
├── tests/
│   └── k6/
│       └── load-test.js
├── go.mod
├── .gitignore
└── README.md
```

## Prerequisites

Before running this project, make sure you have:

- Go installed (the project declares Go 1.26.3 in go.mod)
- PostgreSQL installed and running
- A PostgreSQL database created for the application
- A JWT secret configured

## Environment Configuration

The app loads environment variables from a .env file if present. If the file is not found, it falls back to the system environment variables.

Create a .env file at the project root with values like:

```env
DB_HOST=localhost
DB_USER=postgres
DB_DATABASE=go_social_api
DB_PASSWORD=your_password
PORT=4001
JWT_SECRET=your_super_secure_secret
```

### Environment variables

- DB_HOST: PostgreSQL host
- DB_USER: database username
- DB_DATABASE: database name
- DB_PASSWORD: database password
- PORT: port used by the API server
- JWT_SECRET: secret key used to sign JWT tokens

## Installation

1. Clone the repository:

```bash
git clone https://github.com/your-user/go-social-api.git
cd go-social-api
```

2. Install Go dependencies:

```bash
go mod download
```

3. Create the PostgreSQL database and adjust the database variables in .env.

4. Ensure the database server is running.

## Running the Application

Start the API:

```bash
go run ./cmd/api
```

By default, the application listens on port 4001 unless you set a different PORT value in the environment.

The app also runs automatic migrations for:

- User
- Post
- PostInteraction

This happens in the main entry point before the server starts listening.

## API Endpoints

### 1) Health Check

#### GET /

Returns a basic server status response.

Example response:

```json
{
  "message": "Server running"
}
```

### 2) Register User

#### POST /api/v1/register

Request body:

```json
{
  "username": "john",
  "email": "john@example.com",
  "password": "secret123"
}
```

Response:

```json
{
  "id": 1,
  "username": "john",
  "email": "john@example.com",
  "password": "$2a$10$..."
}
```

Note: the password is stored as a bcrypt hash in the database.

### 3) Login

#### POST /api/v1/login

Request body:

```json
{
  "email": "john@example.com",
  "password": "secret123"
}
```

Successful response:

```json
{
  "user": {
    "id": 1,
    "username": "john",
    "email": "john@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 4) Create a Post

#### POST /api/v1/posts

Requires authentication.

Headers:

```http
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

Request body:

```json
{
  "content": "This is my first post!"
}
```

Response:

```json
{
  "id": 1,
  "user_id": 1,
  "content": "This is my first post!",
  "like_count": 0,
  "deslike_count": 0,
  "interactions": [],
  "created_at": "2026-10-09T00:00:00Z",
  "updated_at": "2026-10-09T00:00:00Z"
}
```

### 5) List Posts

#### GET /api/v1/posts

Requires authentication.

Returns the full post list with interactions.

Example response:

```json
[
  {
    "id": 1,
    "user_id": 1,
    "content": "This is my first post!",
    "like_count": 2,
    "deslike_count": 1,
    "interactions": [
      {
        "id": 1,
        "post_id": 1,
        "user_id": 2,
        "type": "LIKE"
      }
    ]
  }
]
```

### 6) Create or Toggle Interaction

#### POST /api/v1/posts/interation

Requires authentication.

This endpoint accepts a post interaction type. Valid values are:

- LIKE
- DESLIKE

Request body:

```json
{
  "post_id": 1,
  "type": "LIKE"
}
```

Behavior:

- if the user has not interacted before, the interaction is created
- if the same interaction is submitted again, it is removed
- if the user changes from LIKE to DESLIKE or vice versa, the previous count is decremented and the new one is incremented

## Authentication

Protected routes use a JWT in the Authorization header:

```http
Authorization: Bearer <token>
```

The token is generated during login and validated by the auth middleware. If the token is missing, malformed, expired, or invalid, the server returns a 401 Unauthorized response.

## Database Model Notes

The app auto-migrates the following entities:

- User
- Post
- PostInteraction

### User fields

- id
- username
- email
- password
- created_at
- updated_at

### Post fields

- id
- user_id
- content
- like_count
- deslike_count
- created_at
- updated_at

### PostInteraction fields

- id
- post_id
- user_id
- type

## Validation and Error Handling

The API returns descriptive JSON errors for invalid requests, such as:

- missing request body fields
- invalid JSON
- invalid token
- invalid interaction type
- user not found
- email already registered

Common HTTP status codes:

- 200 OK
- 201 Created
- 400 Bad Request
- 401 Unauthorized
- 500 Internal Server Error

## Example Flow

1. Create a user with POST /api/v1/register
2. Log in with POST /api/v1/login
3. Copy the JWT token from the login response
4. Add the token to the Authorization header as Bearer token
5. Create a post with POST /api/v1/posts
6. Retrieve posts with GET /api/v1/posts
7. Interact with a post with POST /api/v1/posts/interation

## Load Test

A sample K6 test is included in:

```text
tests/k6/load-test.js
```

You can run it with:

```bash
k6 run tests/k6/load-test.js
```

This script creates users, logs in, lists posts, creates a post, and performs like/dislike actions.

## Notes

This project is a good example of a small backend service built with Go and the Gin framework. It demonstrates:

- modular package organization
- dependency injection-like composition in main.go
- JWT-based API security
- GORM database integration
- RESTful endpoint patterns for a social application

## License

This project is currently provided without a formal license declaration.

If you want to use it in a production environment, it is recommended to add a license file and improve security, validations, and API documentation.
