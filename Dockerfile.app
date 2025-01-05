# Build-Stage
FROM golang:1.23-alpine AS build
WORKDIR /app

# Copy the source code
COPY . .

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate templ files
RUN templ generate

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/app/main.go
RUN CGO_ENABLED=1 GOOS=linux go build -o cli ./cmd/cli/main.go

# Deploy-Stage
FROM alpine:3.20.2
WORKDIR /app

# Install ca-certificates
RUN apk add --no-cache ca-certificates

# Set environment variable for runtime
ENV GO_ENV=prod

# Copy the binary from the build stage
COPY --from=build /app/main .
COPY --from=build /app/cli .

# Expose the port your application runs on
EXPOSE 8090

# Command to run the application
CMD ["./main"]
