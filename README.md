# Univia Microservices Platform

A comprehensive DevOps SaaS platform built with microservices architecture, providing project management, real-time communication, WebRTC video conferencing, screen sharing, and collaborative tools for modern teams.

## Overview

Univia is a cloud-native platform designed to streamline team collaboration and project management. It combines the power of real-time communication with robust project tracking capabilities, all built on a scalable microservices architecture.

### Key Features

- **Project Management**: Comprehensive project tracking, task management, and team coordination
- **Real-time Communication**: WebRTC-powered video conferencing with HD quality
- **Screen Sharing**: Share your screen with team members during meetings
- **Video Conferencing**: Multi-party video calls with advanced features
- **Meeting Rooms**: Create and manage virtual meeting spaces
- **Team Collaboration**: Real-time collaboration tools and notifications
- **Microservices Architecture**: Scalable, maintainable, and independently deployable services

## Architecture

```
univia/
├── api/              # API Documentation
├── build/            # Build artifacts and Docker images
├── cmd/              # Application entry points
├── configs/          # Configuration files
├── infra/            # Infrastructure as Code (IaC)
├── internal/         # Internal application code
├── pkg/              # Shared packages and libraries
├── scripts/          # Build and deployment scripts
└── test/             # Test suites and fixtures
```

## Prerequisites

- **Go**: 1.21 or higher
- **Docker**: 20.10 or higher
- **Mysql**: 8.0 or higher
- **Redis**: 7.0 or higher
- **Make**: (optional) for build automation

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/deva-labs/univia.git
cd univia
```

### 2. Configure Environment

Copy the example environment file and configure your settings:

```bash
cp configs/.env.example configs/.env
```

### 3. Development Mode

Use the development script to start all services:

```bash
chmod +x run.sh
./run.sh
```

The `run.sh` script will:
- Start required Docker containers (Mysql, Redis)
- Run database migrations
- Seed initial data (roles, permissions, admin user)
- Start the API server with hot reload
- Start the WebRTC signaling server

**Services will be available at:**
- API Gateway: `http://localhost:2000`
- WebRTC Server: `ws://localhost:2112`
- API Documentation: `http://localhost:2000/swagger/index.html`

### 4. Production Deployment

For production deployment, use the release script:

```bash
chmod +x release.sh
./release.sh <alpha/beta> <version> <stage>
# ./release.sh alpha v0.0.2 3 -> v0.0.2-alpha.3
```

The `release.sh` script will:
- Build optimized Docker images
- Run security scans
- Execute test suites
- Create production-ready artifacts
- Generate deployment manifests
- Push images to container registry (if configured)

## Project Structure

### `/api`
Application documentation.
### `/cmd`
Application entry points for different microservices:
- `api`: Main API gateway
- `signaling`: WebRTC signaling server
- `sfu`: SFU server

### `/internal`
Internal application logic organized by domain:
- `api/`: RestAPI business logic
- `signaling/`: WebRTC business logic
- `infrastructure/`: 3rd party client config
- `sfu/`: SFU service business logic

### `/pkg`
Reusable packages shared across services:
- `models/`: Shared models
- `signaling/`: WebRTC implementation
- `types/`: Types public implementation
- `utils/`: Helpers implementation

### `/infra`
Infrastructure as Code (IaC) for deployment: (Not-yet updated)
- `docker/`: Docker and Docker Compose files
- `kubernetes/`: Kubernetes manifests
- `terraform/`: Terraform configurations
- `helm/`: Helm charts

### `/configs`
Configuration files for different environments:
- `.env.development`
- `.env.staging`
- `.env.production`

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./test/integration/...
```

## API Documentation

Once the server is running, access the interactive API documentation:

- **Swagger UI**: `http://localhost:2000/swagger/index.html`

### Authentication

Most endpoints require authentication. To authenticate:

1. Register a new user: `POST /api/v1/auth/register`
2. Login: `POST /api/v1/auth/login`
3. Use the returned token in the `Authorization` header: `Bearer <token>`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [docs/](./docs/)
- **Issues**: [GitHub Issues](https://github.com/yourorg/univia/issues)
- **Email**: support@univia.com
- **Slack**: [Join our community](https://univia.slack.com)

## Roadmap

- [ ] Mobile applications (iOS/Android)
- [ ] Advanced analytics dashboard
- [ ] AI-powered meeting summaries
- [ ] Integration with third-party tools (Slack, Teams, etc.)
- [ ] Recording and playback features
- [ ] Virtual backgrounds and filters
- [ ] Multi-language support

---

**Built with ❤️ by the Univia Team**