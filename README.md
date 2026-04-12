# Posted

A scalable microservices platform for user authentication and content management, built with Go, gRPC, and modern cloud-native practices.

## About

**Posted** is a backend system designed to handle user registration, authentication, session management, and blog post operations. It follows a clean architecture pattern, separating concerns into domain, application, and infrastructure layers. The system ensures high security through mTLS communication between services, JWT-based stateless authentication, and strict input validation using `protovalidate`.

## Features

- **Microservices Architecture**: Decoupled services for Authentication (`auth-service`) and Posts Management (`posts-service`).
- **API Gateway**: A dedicated `http-gateway` exposes RESTful endpoints to the client, handling HTTP-to-gRPC translation.
- **Security First**:
  - **mTLS**: Mutual TLS encryption for service-to-service communication.
  - **JWT**: Secure session management with Access/Refresh token rotation.
  - **Validation**: Strict request validation at the gRPC layer using `buf.build/protovalidate`.
  - **Password Hashing**: Secure storage of user credentials.
- **Data Persistence & Caching**:
  - **PostgreSQL**: Primary relational database with automatic migrations via GORM.
  - **Redis**: High-performance caching for sessions and frequently accessed data.
- **Documentation**: Auto-generated Swagger (OpenAPI) documentation available via the gateway.
- **Containerization**: Full Docker Compose setup for easy local development and deployment.

### Service Breakdown

1.  **Auth Service** (`auth-service`)
    *   Handles user registration, login, profile updates, and deletion.
    *   Manages JWT generation and validation.
    *   Stores user data in PostgreSQL and sessions in Redis.
    *   Exposes gRPC interface defined in `contracts/auth/v1/auth.proto`.

2.  **Posts Service** (`posts-service`)
    *   Manages CRUD operations for blog posts.
    *   Enforces ownership rules (users can only manage their own posts).
    *   Exposes gRPC interface defined in `contracts/posts/v1/posts.proto`.

3.  **HTTP Gateway** (`http-gateway`)
    *   Acts as the single entry point for external clients.
    *   Translates HTTP/JSON requests to gRPC calls.
    *   Provides Swagger UI for API exploration.
    *   Handles cross-cutting concerns like initial request logging.

## Tech Stack

-   **Language**: Go 1.25+
-   **Communication**: gRPC, Protobuf
-   **Web Framework**: Gin (Gateway)
-   **Database**: PostgreSQL (via GORM)
-   **Cache**: Redis (go-redis)
-   **Validation**: `buf.build/protovalidate`, `go-playground/validator`
-   **Security**: mTLS, JWT (golang-jwt), bcrypt
-   **Infrastructure**: Docker, Docker Compose
-   **Docs**: Swaggo

## Getting Started

### Prerequisites

-   Go 1.25+
-   Docker & Docker Compose
-   Make (optional, for convenience commands)
-   `mkcert` (for generating local SSL certificates)

### Quick Start

1.  **Clone the repository**:
    ```bash
    git clone <repository-url>
    cd posted
    ```

2.  **Initialize Environment & Dependencies**:
    Run the `quick` target to generate JWT keys, copy environment variables, and start containers.
    ```bash
    make quick
    ```
    *Note: This creates `.env` from `.env.example` and starts PostgreSQL, Redis, and all services.*

3.  **Access API Documentation**:
    Once the gateway is running, open your browser:
    ```
    http://localhost:8090/swagger/index.html
    ```

### Configuration

Edit the `.env` file to configure database credentials, ports, and feature flags.

-   `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DBNAME`: Database credentials.
-   `APP_CONFIG_PATH`: Path to the specific YAML config for each service.
-   `SERVER_GRPC_MTLS_ENABLE`: Toggle mTLS (true/false).

## Project Structure

```
.
├── auth-service/          # Authentication microservice
├── posts-service/         # Posts management microservice
├── http-gateway/          # REST API Gateway
├── contracts/             # Shared Protobuf definitions
│   ├── auth/v1/
│   └── posts/v1/
├── docker-compose.yml     # Orchestration
└── Makefile               # Build automation
```

## API Endpoints (via Gateway)

### Authentication
-   `POST /api/v1/register` - Register a new user.
-   `POST /api/v1/login` - Login and receive tokens.
-   `GET /api/v1/user` - Get current user profile (Protected).
-   `PATCH /api/v1/user` - Update user profile (Protected).
-   `DELETE /api/v1/user` - Delete account (Protected).

### Posts
-   `POST /api/v1/post` - Create a new post (Protected).
-   `PATCH /api/v1/posts/:id` - Update a post (Protected).
-   `GET /api/v1/post/:id` - Get a specific post.
-   `GET /api/v1/posts` - Get current user's posts (Protected).
-   `DELETE /api/v1/post/:id` - Delete a post (Protected).

*Protected endpoints require an `Authorization: <token>` header.*

## License

MIT License - see [LICENSE](LICENSE) for details.
