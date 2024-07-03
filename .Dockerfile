# Start with the official Golang base image
FROM golang:1.20-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o /unified-go-backend cmd/main.go

# Start a new stage from scratch
FROM alpine:latest

# Set environment variables
ENV MONGO_URI=mongodb://mongo:27017/mdmdb
ENV REDIS_ADDR=redis:6379
ENV SMTP_HOST=smtp.mailtrap.io
ENV SMTP_PORT=587
ENV SMTP_USER=your_mailtrap_user
ENV SMTP_PASSWORD=your_mailtrap_password
ENV JWT_SECRET=your_jwt_secret
ENV ELASTICSEARCH_URL=http://elasticsearch:9200

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /unified-go-backend /unified-go-backend

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["/unified-go-backend"]
