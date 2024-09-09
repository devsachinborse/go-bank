# Gobank

Gobank is a simple banking application that uses PostgreSQL for storage and provides various API endpoints for managing bank accounts. This README will guide you through setting up and running the application.

## Prerequisites

- Docker: [Installation Guide](https://docs.docker.com/get-docker/)
- Docker Compose (optional, for multi-container setups)

## Running the Application

### Starting PostgreSQL with Docker

To run PostgreSQL as a Docker container, use the following command:

```sh
docker run --name gobank -e POSTGRES_PASSWORD=gobank -p 5432:5432 -d postgres:latest
```
This command will:

Create a new container named gobank
Set the PostgreSQL superuser password to gobank
Expose port 5432 on the host machine
Application Setup

### Clone the Repository

```
https://github.com/devsachinborse/go-bank.git
```
### Build and Run the Application

Ensure you have Go installed on your machine. Then, navigate to the project directory and build the application:
Run the application:
```
make run
```
By default, the application will start listening on port 3000.


