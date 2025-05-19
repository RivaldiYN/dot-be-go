# Go Book API

## Arsitektur

Proyek ini menggunakan **Clean Architecture** dengan struktur direktori yang memisahkan tanggung jawab ke dalam beberapa lapisan. Tujuannya untuk:

- Memudahkan pengujian unit di tiap lapisan
- Memungkinkan fleksibilitas penggantian teknologi/framework
- Meningkatkan skalabilitas dan maintainability

### Struktur Direktori

```text
├── cmd/
│   └── main.go              # Entry point aplikasi
│
├── config/
│   └── config.go            # Konfigurasi database dan environment
│
├── internal/
│   ├── app/
│   │   └── api/
│   │       ├── handlers/            # HTTP handlers (controller layer)
│   │       │   ├── auth_handler.go
│   │       │   ├── book_handler.go
│   │       │   ├── category_handler.go
│   │       │   └── handler.go
│   │       ├── middleware/          # Middleware (auth, logger, dll)
│   │       │   └── jwt_middleware.go
│   │       └── routes/              # HTTP routes definition
│   │           └── routes.go
│   │
│   ├── domain/
│   │   ├── entity/                  # Entitas domain (data structure)
│   │   │   ├── book.go
│   │   │   ├── category.go
│   │   │   └── user.go
│   │   ├── repository/             # Abstraksi akses data (interface & impl)
│   │   │   ├── book_repository.go
│   │   │   ├── category_repository.go
│   │   │   └── user_repository.go
│   │   └── service/                # Business logic layer
│   │       ├── auth_service.go
│   │       ├── book_service.go
│   │       └── category_service.go
│
├── pkg/                            # Shared utilities
│   ├── hash/
│   │   └── hash.go
│   └── jwt/
│       └── jwt.go
│
├── test/
│   └── e2e/
│       └── auth_test.go            # End-to-end tests
│
├── go.mod
└── go.sum
```

## Teknologi yang Digunakan

- [Echo](https://echo.labstack.com/) – Web framework ringan dan cepat
- [GORM](https://gorm.io/) – ORM untuk database
- [PostgreSQL](https://www.postgresql.org/) – Relational Database
- [JWT](https://jwt.io/) – Authentication
- `testing` + `httptest` – End-to-End dan Unit Testing


## Menjalankan Aplikasi & Testing

1. Clone repository
2. Jalankan perintah:

```bash
go run cmd/main.go & go test -v ./test/e2e