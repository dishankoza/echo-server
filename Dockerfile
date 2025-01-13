# Use official Golang image as base
FROM golang:1.23.4 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the project
COPY . .

# Build the Go app (checking the output binary)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Debug step to confirm the binary is built correctly
RUN ls -l /app

# Start a new stage from a smaller image to run the binary
FROM alpine:latest  

# Install necessary dependencies for your app to run
RUN apk --no-cache add ca-certificates

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/main

# Debug step to confirm binary was copied correctly
RUN ls -l /usr/local/bin/main

# Make the binary executable
RUN chmod +x /usr/local/bin/main

# Expose the port
EXPOSE 7379

# Command to run the executable
CMD ["/usr/local/bin/main"]