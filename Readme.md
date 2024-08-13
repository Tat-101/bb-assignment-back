[![codecov](https://codecov.io/github/Tat-101/bb-assignment-back/graph/badge.svg?token=A6A0INM06I)](https://codecov.io/github/Tat-101/bb-assignment-back)

# BB Assignment - Backend Application

This repository contains the backend application for the "BB Assignment". This project is built using Go and provides the necessary API endpoints for managing users in the system.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [Seeding the Database](#seeding-the-database)
- [Testing](#testing)
- [Documentation](#documentation)

## Overview

The "BB Assignment" backend application is part of a user management system. It provides RESTful API endpoints for managing users, including role-based access control. The backend is responsible for handling business logic, authentication, and data management.

## Features

- **User Management**: Create, update, delete, and list users.
- **Role-Based Access Control**: Only administrators can perform certain actions, such as deleting users.
- **Authentication**: Secure login with session management.
- **Database Seeding**: Seed initial data, including an admin user.

## Installation

To install the backend application, follow these steps:

```bash
git clone https://github.com/Tat-101/bb-assignment-back.git
cd bb-assignment-back
go mod tidy
cp .env.example .env
```

This command will install all the necessary dependencies.

## Running the Application

Before running the backend application, ensure that PostgreSQL is running. You can either start PostgreSQL manually or use Docker Compose to start all necessary services, including PostgreSQL and the backend:

- **Option 1: Ensure PostgreSQL is running manually**:  
  Start your PostgreSQL service using your preferred method (e.g., systemctl, pg_ctl, etc.). Once PostgreSQL is running, you can start the backend application with the following command:

  ```bash
    go run .
  ```

* **Option 2: Using Docker Compose**:
  ```bash
    docker-compose up -d
  ```

The backend will start and be accessible at the specified host and port in your configuration.

## Seeding the Database

Before using the application, you need to seed the database with initial data, including an admin user. Run the following command:

```bash
go run tools/seed/main.go
```

This will create an admin user with the following credentials:

- **Username**: `admin@bb.com`
- **Password**: `123456`

Only the admin user can perform delete operations on other users.

## Testing

To run the tests included in the project, use the following command:

```bash
go test ./...
```

This will execute all tests across the project.

## Documentation

For detailed API documentation, visit the [Postman Documentation](https://documenter.getpostman.com/view/1837888/2sA3s4mWNC).
