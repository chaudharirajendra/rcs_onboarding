# RCS Onboarding Backend

## Overview
Golang backend for RCS form management with versioning, validation, workflows, audits.

## Setup
- Go 1.22+, MySQL.
- Env: DB_DSN, JWT_SECRET.
- `go mod tidy`
- `go run cmd/main.go`

## Docker
`docker build -t rcs-onboarding .`
`docker run -p 8080:8080 -e DB_DSN=... rcs-onboarding`

## API Docs
Run `swag init` for swagger.json (requires github.com/swaggo/swag). Access /swagger/index.html (add gin-swagger middleware).

## Schemas
Qualification: 20 vetting fields. Customer Order: 13 setup fields. See utils/validator.go for seeding.

## Testing
Postman: Login {"username":"customer","password":"password"}, submit forms, review.