write a go boiler plate in this structure
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── config
│   │   ├── config.go
│   │   └── database.go
│   ├── handlers
│   │   ├── auth_handler.go
│   │   ├── error_handler.go
│   │   └── user_handler.go
│   ├── middleware
│   │   ├── auth.go
│   │   ├── cache.go
│   │   ├── rate_limit.go
│   │   └── validator.go
│   ├── models
│   │   └── user.go
│   ├── router
│   │   ├── admin_routes.go
│   │   ├── auth_routes.go
│   │   ├── health_routes.go
│   │   ├── router.go
│   │   └── user_routes.go
│   └── services
│       └── user_service.go
├── go.mod
└── go.sum

Database to use mongodb and caching layer as redis, boiler plate handles user registration and login using mobile number.
with
go version 1.23
go-redis/v9