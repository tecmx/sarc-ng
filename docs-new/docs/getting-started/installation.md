---
sidebar_position: 1
---

# Installation

This guide covers the installation and setup process for SARC-NG.

## Prerequisites

Before installing SARC-NG, ensure you have the following prerequisites:

- **Go** (version 1.24 or later)
- **Docker & Docker Compose** (for containerized development)
- **Make** (for build automation)
- **Git** (for version control)

### Installing Prerequisites

#### Linux (Ubuntu/Debian)

```bash
# Install Go
wget https://go.dev/dl/go1.24.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile

# Install Docker
sudo apt-get update
sudo apt-get install -y docker.io docker-compose
sudo usermod -aG docker $USER

# Install Make & Git
sudo apt-get install -y make git
```

#### macOS

```bash
# Install Homebrew if not already installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install prerequisites
brew install go@1.24
brew install docker docker-compose
brew install make git
```

#### Windows

1. **Install Go**: Download and run the installer from the [Go downloads page](https://golang.org/dl/).

2. **Install Docker Desktop**: Download and install from the [Docker website](https://www.docker.com/products/docker-desktop).

3. **Install Git**: Download and install from the [Git website](https://git-scm.com/download/win).

4. **Install Make**:
   - Install [Chocolatey](https://chocolatey.org/install)
   - Run: `choco install make`

## Installing SARC-NG

### Clone the Repository

```bash
# Clone the repository
git clone <repository-url>
cd sarc-ng
```

### Setup Development Environment

```bash
# Install dependencies and tools
make setup
```

This will:
- Download Go dependencies
- Install development tools
- Initialize the project

### Verify Installation

To verify that everything is installed correctly:

```bash
# Run the tests
make test

# Build the application
make build
```

You should see successful test execution and binary creation without errors.

## Configuration

Basic configuration is handled through environment variables and configuration files:

### Environment Variables

```bash
# Database configuration
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=password
export DB_NAME=sarcng

# Server configuration
export PORT=8080
export ENVIRONMENT=development
```

### Configuration Files

The project includes default configurations:

- `configs/default.yaml` - Base configuration
- `configs/development.yaml` - Development overrides

## Next Steps

Now that you have installed SARC-NG, you can:

- [Learn about development workflows](development)
- [Set up Docker environment](docker)
- [Follow the quick start guide](quick-start)
