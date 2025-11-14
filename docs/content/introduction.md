---
sidebar_position: 1
---

# Introduction

**SARC-NG** - A modern resource management and scheduling API for educational institutions.

## What is SARC-NG?

An API-driven system for managing:
- **Buildings** - Campus facilities and locations
- **Classrooms** - Learning spaces and labs
- **Resources** - Equipment and facilities
- **Scheduling** - Lessons and reservations
- **Classes** - Course management with instructors

## Key Features

- 🏢 **Building & Resource Management** - Track buildings, rooms, equipment with custom attributes
- 📅 **Smart Scheduling** - Flexible repeat patterns (daily, weekly, biweekly) with conflict detection
- 👥 **Class Management** - Organize classes with instructors, semesters, and enrollment
- 🔐 **Secure Authentication** - JWT-based with bearer tokens and role-based access
- 📊 **RESTful API** - OpenAPI 3.1 spec with standardized responses

## Quick Start

```bash
git clone https://github.com/tecmx/sarc-ng.git
cd sarc-ng
make docker-up
```

**API**: http://localhost:8080/api/v1
**Swagger**: http://localhost:8080/swagger/index.html

## API Structure

Standard REST operations for all entities:

```
GET    /api/v1/{entity}        # List all
POST   /api/v1/{entity}        # Create
GET    /api/v1/{entity}/:id    # Get by ID
PUT    /api/v1/{entity}/:id    # Update
DELETE /api/v1/{entity}/:id    # Delete
```

**Entities**: `buildings`, `classes`, `lessons`, `resources`, `reservations`

## Next Steps

- 📚 [Getting Started Guide](getting-started.md)
- 🚀 [Development Guide](development.md)
- 🏗️ [Architecture Overview](architecture.md)
- 🌐 [Deployment Guide](deployment.md)
