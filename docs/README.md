# Go Vue Admin Documentation

## Table of Contents

1. [Introduction](./introduction.md)
2. [Getting Started](./getting-started.md)
3. [Architecture](./architecture.md)
4. [API Documentation](./api/README.md)
5. [Frontend Documentation](./frontend/README.md)
6. [Database Schema](./database/README.md)
7. [Deployment Guide](./deployment/README.md)

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Node.js 18 or higher
- Docker (optional)
- Make (optional)

### Development Setup
1. Clone the repository
2. Run `make frontend-install` to install frontend dependencies
3. Run `make dev` to start both frontend and backend in development mode

### Production Build
1. Run `make build` to build the backend
2. Run `make frontend-build` to build the frontend
3. Run `make docker-build` to create a Docker image

For detailed documentation, please refer to the respective sections above.

## Contributing

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details. 