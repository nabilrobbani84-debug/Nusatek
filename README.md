# Nusatek Property Management Backend Service

This project is a high-performance backend service designed to demonstrate a scalable architecture for a Property Management System. It is engineered to meet the specific technical requirements of **PT Nusantara Sukses Teknologi (Nusatek.id)**, showcasing expertise in Golang, Clean Architecture, and Distributed Systems.

## 🚀 Key Features & Tech Stack

This project maps directly to the job requirements:

*   **Language:** Golang (1.21+)
*   **Architecture:** **Clean Architecture** (Domain-Driven Design principles) to ensure maintainability and testability.
*   **Database:** **PostgreSQL** for reliable relational data storage.
*   **Performance Optimization:** **Redis** for caching hot data (property listings).
*   **Asynchronous Processing:** **RabbitMQ** for decoupling heavy tasks (e.g., email notifications, audit logs) from the main API response loop.
*   **Containerization:** **Docker & Docker Compose** for easy deployment and environment consistency.
*   **Documentation:** Swagger/OpenAPI (planned).

## 📂 Project Structure (Standard Go Layout)

```
nusatek-property-backend/
├── cmd/
│   └── api/            # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── delivery/       # HTTP Handlers (Gin/Echo)
│   ├── domain/         # Business logic interfaces & entities (The Core)
│   ├── repository/     # Database implementations (Postgres/Redis)
│   └── usecase/        # Application business logic
├── pkg/
│   ├── database/       # DB connection helpers
│   ├── logger/         # Structured logging
│   └── rabbitmq/       # Message queue helpers
└── docker-compose.yml  # Infrastructure setup
```

## 🛠️ How to Run

1.  **Prerequisites:** Docker and Go installed.
2.  **Start Infrastructure:**
    ```bash
    docker-compose up -d
    ```
3.  **Run Application:**
    ```bash
    go run cmd/api/main.go
    ```
