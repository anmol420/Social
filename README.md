# Social API

A RESTful API for a social media platform built with Go, featuring user authentication, posts, comments, and social interactions.

## Features

- **User Management**
  - User registration with email verification
  - JWT-based authentication
  - User activation via token
  - User profiles and feeds

- **Social Interactions**
  - Create, read, update, and delete posts
  - Comment on posts
  - Follow/unfollow users
  - Personalized user feed

- **Security**
  - JWT token authentication
  - Basic auth for admin endpoints
  - Role-based access control (admin, moderator)
  - Password hashing
  - Rate limiting capabilities

- **Monitoring & Health**
  - Health check endpoint
  - Debug metrics via expvar
  - Comprehensive logging

## Tech Stack

- **Language**: Go
- **Router**: Chi
- **Database**: PostgreSQL (assumed from structure)
- **Authentication**: JWT tokens
- **Email**: Amazon SES
- **Cache**: Redis

## Project Structure

```
Social/
├── .github/          # GitHub workflows and CI/CD
├── bin/             # Compiled binaries
├── cmd/             # Application entrypoints
│   ├── api/         # Main API server
│   └── migrate/     # Database migrations
├── docs/            # API documentation
├── internal/        # Private application code
│   ├── auth/        # Authentication logic
│   ├── db/          # Database layer
│   ├── env/         # Environment configuration
│   ├── mailer/      # Email service
│   ├── ratelimiter/ # Rate limiting
│   └── store/       # Data access layer
├── scripts/         # Build and deployment scripts
├── web/             # Web assets (if any)
├── .env.sample      # Environment variables template
├── docker-compose.yml # Docker setup
├── go.mod           # Go dependencies
├── go.sum           # Dependency checksums
└── Makefile         # Build automation
```

## API Endpoints

### Health & Monitoring

```
GET  /v1/health              # Health check (public)
GET  /v1/debug/vars          # Runtime metrics (basic auth required)
```

### Authentication

```
POST /v1/auth/register       # Register new user
POST /v1/auth/token          # Generate JWT token
```

### Users

```
PUT  /v1/users/activate/{token}     # Activate user account
GET  /v1/users/{userID}             # Get user profile (auth required)
POST /v1/users/{userID}/follow      # Follow user (auth required)
POST /v1/users/{userID}/unfollow    # Unfollow user (auth required)
GET  /v1/users/feed                 # Get user feed (auth required)
```

### Posts

```
POST   /v1/posts/create              # Create post (auth required)
GET    /v1/posts/{postID}            # Get post by ID (auth required)
PATCH  /v1/posts/{postID}            # Update post (moderator+ required)
DELETE /v1/posts/{postID}            # Delete post (admin required)
POST   /v1/posts/{postID}/comments   # Create comment (auth required)
```

## Getting Started

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 14 or higher
- Docker
- Make

### Installation

1. Clone the repository:
```bash
git clone https://github.com/anmol420/social.git
cd social
```

2. Copy environment variables:
```bash
cp .env.sample .env
```

3. Configure your `.env` file with appropriate values:
```env
ADDR=""
EXTERNAL_URL=""
DATABASE_ADDR=""
MAX_OPEN_CONNS=""
MAX_IDLE_CONNS=""
MAX_IDLE_TIME=""
MAILER_FROM_EMAIL=""
MAILER_REGION=""
FRONTEND_URL=""
AUTH_BASIC_USERNAME=""
AUTH_BASIC_PASSWORD=""
AUTH_JWT_TOKEN_SECRET=""
REDIS_ADDR=""
REDIS_DB=""
REDIS_ENABLED=""
RATELIMITER_REQUESTS_COUNT=""
RATELIMITER_ENABLED=""
```

4. Install dependencies:
```bash
go mod download
```

5. Run database migrations:
```bash
make migrate-up
```

6. Start the server:
```bash
air
# or
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080` (or your configured port).

## Docker Deployment

### Using Docker Compose

1. Build and start services:
```bash
docker-compose up -d
```

2. Run migrations:
```bash
docker-compose exec api make migrate-up
```

3. View logs:
```bash
docker-compose logs -f api
```

4. Stop services:
```bash
docker-compose down
```

## Authentication

### Register a User

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Get Authentication Token

```bash
curl -X POST http://localhost:8080/v1/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Using the Token

Include the JWT token in the Authorization header:

```bash
curl -X GET http://localhost:8080/v1/users/feed \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Role-Based Access Control

The API implements role-based access control with the following roles:

- **User**: Basic access to create posts, comments, follow users
- **Moderator**: Can update any post
- **Admin**: Full access including post deletion

## Rate Limiting

Rate limiting is implemented to prevent abuse. Default limits:
- 200 requests per minute per IP
- Configurable per endpoint

## Error Handling

The API returns consistent error responses:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {}
  }
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Database Migrations

### Create a New Migration

```bash
make migrate-create <migration_name>
```

### Apply Migrations

```bash
# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down
```

## Monitoring

### Health Check

```bash
curl http://localhost:8080/v1/health
```

Response:
```json
{
  "status": "ok",
  "timestamp": "2024-02-07T10:30:00Z"
}
```

### Debug Metrics

Access runtime metrics (requires basic auth):

```bash
curl -u admin:password http://localhost:8080/v1/debug/vars
```

## Security Considerations

- All passwords are hashed using bcrypt
- JWT tokens expire after a configured duration
- HTTPS should be enforced in production
- SQL injection prevention through parameterized queries
- Input validation on all endpoints
- CORS configured for specific origins
- Rate limiting to prevent abuse
- Secure headers configured

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Coding Standards

- Follow Go best practices and idioms
- Write tests for new features
- Update documentation for API changes
- Use meaningful commit messages
- Keep functions focused and small
- Add comments for complex logic

## Troubleshooting

### Common Issues

**Database Connection Failed**
- Verify database credentials in `.env`
- Ensure PostgreSQL is running
- Check network connectivity

**JWT Token Invalid**
- Verify `JWT_SECRET` is set
- Check token expiration
- Ensure proper Authorization header format

**Migration Errors**
- Check migration files for syntax errors
- Verify database permissions
- Review migration logs

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes.
