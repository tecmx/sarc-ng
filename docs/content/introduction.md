---
sidebar_position: 1
tags:
  - getting-started
  - introduction
---

# Introduction

Welcome to **SARC-NG** (Resource Management and Scheduling System) - a comprehensive solution for managing educational resources and scheduling.

## What is SARC-NG?

SARC-NG is a modern resource management and scheduling system designed for educational institutions. It provides a complete API-driven solution for:

- **Building Management** - Organize and track your campus buildings
- **Classroom Management** - Manage classrooms, labs, and other learning spaces
- **Resource Management** - Handle equipment, facilities, and other resources
- **Scheduling** - Create and manage lessons, reservations, and time slots
- **User Management** - Authentication and authorization system

## Key Features

### 🏢 Building & Resource Management
- Track buildings with detailed information (floors, addresses, codes)
- Manage different types of resources (classrooms, equipment, labs)
- Flexible resource details with custom attributes
- Capacity tracking and location management

### 📅 Smart Scheduling
- Create lessons with flexible repeat patterns (once, daily, weekly, biweekly)
- Resource reservations with conflict detection
- Time-based filtering and availability checking
- Automated scheduling validation

### 👥 Class Management
- Organize classes with instructor information
- Track semester information and capacity
- Link classes to scheduled lessons
- Student enrollment capabilities

### 🔐 Authentication & Security
- JWT-based authentication system
- Secure API endpoints with bearer token authentication
- Role-based access control
- User session management

## API Architecture

SARC-NG follows REST API principles with:

- **OpenAPI 3.1 specification** for complete API documentation
- **JSON-based** request/response format
- **HTTP status codes** for proper error handling
- **Standardized error responses** across all endpoints
- **Bearer token authentication** for security

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose
- PostgreSQL database
- Redis (for caching and sessions)

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/tecmx/sarc-ng.git
   cd sarc-ng
   ```

2. **Start with Docker Compose**
   ```bash
   cd docker
   docker compose up -d
   ```

3. **API will be available at**
   ```
   http://localhost:8080/v1
   ```

### Explore the API

Check out our [**API Reference**](/content/category/api-reference) for complete endpoint documentation with:

- Interactive API explorer
- Request/response examples
- Authentication details
- Error code references

## Next Steps

- 📚 Browse the [API Reference](/content/category/api-reference) to explore all available endpoints
- 🚀 Try the interactive API examples to test functionality
- 💡 Check out our tutorials for common use cases
- 🐛 Report issues on our [GitHub repository](https://github.com/tecmx/sarc-ng)

Ready to get started? Jump into our [API Reference](/content/category/api-reference) and start building!
