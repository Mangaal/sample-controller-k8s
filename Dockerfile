# Stage 1: Build the Go application
FROM golang:1.13 AS build

# Set the working directory within the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire application source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd

# Stage 2: Create a minimal production image
FROM alpine:13

# Set the working directory within the container
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=build /app/myapp .

# Expose a port if your application listens on one
# EXPOSE 8080

# Define the command to run your application
CMD ["./myapp"]
