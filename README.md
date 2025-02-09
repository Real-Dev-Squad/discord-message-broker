# Discord Message Broker Setup and Running Instructions

This document provides instructions on how to set up and run the Go project using the provided `Makefile` commands.

## Requirements

- Go 1.18+
- Make
- Air (for `make air`)

Ensure you have these dependencies installed before running the commands.

## Prerequisites

Before running the project, ensure that you have the following installed:

- **Go 1.18+**: The Go programming language (required to run and build the project).
- **Air**: A live-reloading tool for Go that will automatically restart the project on file changes.
- **Make**: A build automation tool used to manage tasks defined in the `Makefile`.

## Installation

1. **Install Go**
   If you don't have Go installed, follow the official guide to install it:[Go Installation Guide](https://go.dev/doc/install).
2. **Install Make**
   To install Make, follow the installations steps from here based on your OS:
   [Make Installation Guide.](https://www.geeksforgeeks.org/how-to-install-make-on-ubuntu/)
3. **Install Air**
   To install Air, follow the installation steps here:
   [Air Installation Guide](https://github.com/air-verse/air)

## Running Discord Message Broker

### Running both RabbitMQ & Consumer using Docker

Set the following env var(s)

```
QUEUE_URL = "amqp://rabbitmq:5672"
DISCORD_QUEUE = <ANY_NAME> #Default: "DISCORD_QUEUE"
```

Run the compose command:

```sh
docker-compose up --build
```

### Running only RabbitMQ with Docker

Set the following env var(s)

```
QUEUE_URL = "amqp://localhost:5672"
DISCORD_QUEUE = <ANY_NAME> #Default: "DISCORD_QUEUE"
```

```bash
docker compose -f 'docker-compose.yml' up -d --build 'rabbitmq'
```

> [!IMPORTANT]
> To check if the queue is running or not, visit `http://localhost:5672`

### Running Consumer Manually

1. **Install Packages**

```bash
   go mod download
```

2. **Verify Packages**
   If it's your first time running the project, ensure all dependencies are set up:

   ```bash
   go mod tidy
   ```

3. **Running the Project**

   ```bash
   go run .
   ```

4. **Running the Project Using Air**

   ```bash
   air
   ```

## Running the Project Using Make

You can run the project using the `Makefile`, which provides several commands for various tasks. Below are the steps to run the project:

1. **Install Packages**

   ```bash
   make download
   ```

2. **Verify Packages**
   If it's your first time running the project, ensure all dependencies are set up:

   ```bash
   make tidy
   ```

3. **Running the Project**

   ```bash
   make run
   ```

4. **Running the Project Using Air**

   ```bash
   make air
   ```

## Other Commands In Usage

1. **To run tests**:

   ```bash
   make test #or go list ./... | grep -v "/config$$" | grep -v "/routes$$" | xargs go test -v
   ```

2. **To generate a coverage report**:

   ```bash
   make coverage #or go list ./... | grep -v "/config$$" | grep -v "/routes$$" | xargs go test -v -coverprofile=coverage.out
   ```

3. **To automatically re-run the application on changes**:

   ```bash
   make air #or air
   ```

4. **To clean up the generated files**:

   ```bash
   make clean #or rm -rf coverage coverage.out coverage.html
   ```
