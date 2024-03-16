# Stage 1: Build the application
FROM golang:1.22 AS build

WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Run tests
RUN go test ./...

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /receipt-processor

# Stage 2: Create a minimal image with only the binary
FROM alpine:latest

# Copy the built binary from the previous stage
COPY --from=build /receipt-processor /receipt-processor
COPY .env /

EXPOSE ${PORT}

CMD ["/receipt-processor"]