# Go REST API

This is a REST API built with Go, Gin, MongoDB, Redis, and Swagger for API documentation. The API includes user registration, email verification, user login, and profile management functionalities.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Running the Project](#running-the-project)
- [API Documentation](#api-documentation)
- [Deployment](#deployment)
- [Author](#author)

## Features

- User Registration with Email Verification
- User Login with JWT Authentication
- User Profile Management
- Email Verification Job Worker
- Logging with Logrus and Elasticsearch
- API Documentation with Swagger

## Prerequisites

- Go 1.17 or later
- MongoDB
- Redis
- Docker (optional, for deployment)

## Setup

1. **Clone the Repository**

    ```sh
    git clone https://github.com/yourusername/go-rest-api.git
    cd go-rest-api
    ```

2. **Install Dependencies**

    ```sh
    go mod tidy
    ```

3. **Environment Variables**

    Create a `.env` file in the root directory of the project and add the following environment variables:

    ```env
    MONGO_URI=mongodb://localhost:27017
    JWT_SECRET=your_secret_key
    REDIS_ADDR=localhost:6379
    REDIS_PASSWORD=
    SMTP_HOST=smtp.example.com
    SMTP_PORT=587
    SMTP_USER=your_email@example.com
    SMTP_PASSWORD=your_email_password
    ELASTICSEARCH_URL=http://localhost:9200
    ```

4. **Generate Swagger Documentation**

    Install `swag` and generate the Swagger documentation:

    ```sh
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init
    ```

## Running the Project

1. **Run the Application**

    ```sh
    go run cmd/main.go
    ```

2. **Access API**

    The API will be available at `http://localhost:8080`.

3. **Access Swagger Documentation**

    Open your browser and navigate to `http://localhost:8080/swagger/index.html` to access the API documentation.

## API Documentation

The API documentation is generated using Swagger. You can access it at `http://localhost:8080/swagger/index.html` once the application is running.

## Deployment

You can deploy the application using Docker. Here are the steps:

1. **Create a Dockerfile**

    Create a `Dockerfile` in the root directory of the project:

    ```dockerfile
    # Start with a base image containing Go
    FROM golang:1.17-alpine

    # Set the Current Working Directory inside the container
    WORKDIR /app

    # Copy go mod and sum files
    COPY go.mod go.sum ./

    # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
    RUN go mod download

    # Copy the source from the current directory to the Working Directory inside the container
    COPY . .

    # Build the Go app
    RUN go build -o main .

    # Expose port 8080 to the outside world
    EXPOSE 8080

    # Command to run the executable
    CMD ["./main"]
    ```

2. **Build the Docker Image**

    ```sh
    docker build -t go-rest-api .
    ```

3. **Run the Docker Container**

    ```sh
    docker run -d -p 8080:8080 --name go-rest-api --env-file .env go-rest-api
    ```

    Ensure that MongoDB, Redis, and Elasticsearch are accessible to the container, either by running them as separate containers or ensuring they are accessible on the host network.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Author

- **Mohammad Amin Najafi** - *Initial work* - [manvfx](https://github.com/manvfx)
