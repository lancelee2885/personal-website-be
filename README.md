# Personal Website Backend

This is the backend application for a personal website built with Go. It provides a RESTful API for managing projects and jobs.

## Features

- CRUD operations for projects and jobs
- Dockerized application and PostgreSQL database
- Database migrations using golang-migrate
- Logging with zerolog
- Gin web framework for handling HTTP requests
- GORM as the ORM library for database interactions

## Prerequisites

- Go 1.18 or higher
- Docker
- Docker Compose

## Getting Started

1. Build and run the application using Docker Compose:
    
    ```bash
    docker-compose up --build
    ```

2. Access the API at `http://localhost:8080`.

## API Endpoints

- `POST /entities`: Create a new entity (project or job)
- `GET /entities?tableName=<table>`: List all entities from the specified table
- `GET /entity/:id?type=<table>`: Get an entity by ID from the specified table
- `PUT /entities/:id`: Update an entity by ID
- `DELETE /entities/:id`: Delete an entity by ID
- `PATCH /entities/:id/archive`: Archive an entity by ID

Replace `<table>` with either `projects` or `jobs` to specify the desired table.

## Configuration

The application can be configured using environment variables. The following variables are available:

- `DATABASE_HOST`: Database host (default: `db`)
- `DATABASE_PORT`: Database port (default: `5432`)
- `DATABASE_USER`: Database username (default: `myuser`)
- `DATABASE_PASSWORD`: Database password (default: `mypassword`)
- `DATABASE_NAME`: Database name (default: `mydb`)

You can modify these variables in the `docker-compose.yml` file or set them as environment variables when running the application.

## Database Migrations

The database migrations are managed using the golang-migrate library. The migration files are located in the `migrations` directory.

To create a new migration, run the following command:
    
    ```bash
    migrate create -ext sql -dir migrations -seq <migration_name>
    ```

## TODO

- Add authentication and authorization
- Add tests
- Add CI/CD pipeline

