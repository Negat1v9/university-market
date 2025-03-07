# Student Work Marketplace

A REST API service built with Go that provides a marketplace platform for students to post and find academic work opportunities. The platform serves as a bridge between students who need academic assistance and those who can provide it.

## Overview

This project implements a marketplace where:

- Users can post academic tasks with detailed requirements
- Workers (other students) can browse, respond to, and complete tasks
- Secure payment handling through Telegram payments
- Real-time notifications via Telegram Bot
- Rating and feedback system for workers
- File sharing capabilities for task materials

## Features

- **Authentication & Authorization**

  - JWT-based authentication
  - Role-based access control (regular users and workers)
  - Telegram Mini App authentication integration

- **Task Management**

  - Create, update, and delete tasks
  - Task status tracking
  - File attachments support
  - Task promotion system

- **Worker System**

  - Worker profile management
  - Response to tasks
  - Rating and review system
  - Task completion verification

- **Payment System**

  - Secure payment processing via Telegram
  - Balance management
  - Payment history tracking

- **Communication**

  - Integrated Telegram bot notifications
  - Direct messaging system
  - File sharing capabilities

- **Reporting System**
  - User reporting functionality
  - Worker performance tracking
  - Task completion statistics

## Technology Stack

### Core

- Go 1.21+
- MongoDB for data storage
- Telegram Bot API
- JSON Web Tokens (JWT)

### Main Dependencies

- `go-telegram-bot-api/telegram-bot-api/v5`: Telegram Bot integration
- `golang-jwt/jwt/v5`: JWT authentication
- `mongodb/mongo-go-driver`: MongoDB driver
- `stretchr/testify`: Testing framework
- `uber-go/mock`: Mocking for tests
- `log/slog`: Structured logging

## Project Structure

```
├── cmd/
│   └── app/                 # Application entry point
├── internal/
│   ├── app/                # Application initialization
│   ├── config/             # Configuration management
│   ├── services/           # Business logic layer
│   ├── storage/            # Data access layer
│   ├── tgBot/             # Telegram bot implementation
│   └── web/               # HTTP handlers and routing
├── model/                  # Data models
└── pkg/                    # Shared packages
```

## Key Components

### Services

- **Auth Service**: Handles user authentication and authorization
- **Task Service**: Manages task creation, updates, and lifecycle
- **Worker Service**: Handles worker registration and management
- **Payment Service**: Processes payments and manages balances
- **Comment Service**: Manages feedback and ratings
- **Report Service**: Handles user reports and issues

### Storage

- MongoDB repositories for all entities
- Transaction support for critical operations
- Indexed collections for optimal performance

### API Endpoints

The API provides endpoints for:

- User management
- Task operations
- Worker operations
- Payment processing
- Comments and ratings
- Reports and feedback

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MongoDB
- Telegram Bot Token
- Environment variables configuration

### Environment Variables

```
ENV=prod|local
MONGO_URL=mongodb://localhost:27017
MONGO_DB_NAME=marketplace
TG_BOT_TOKEN=your_bot_token
NUMBER_WORKER=10
WEB_APP_BASE_URL=https://your-domain
SERVER_PORT=:8080
JWT_SECRET=your_secret
CONTEXT_TIMEOUT=5
ADMINS_IDS=id1,id2
```

### Running the Application

```bash
# Clone the repository
git clone [repository-url]

# Install dependencies
go mod download

# Run the application
docker compose up --build -d
```

## Testing

The project includes comprehensive unit tests for all major components:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Contributing

Please read CONTRIBUTING.md for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

For any questions or suggestions, please create an issue in the repository.

Note: This project is designed for educational purposes and should be used in accordance with academic integrity policies.
