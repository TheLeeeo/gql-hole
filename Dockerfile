FROM golang:1.21 as builder

LABEL org.opencontainers.image.source=https://github.com/TheLeeeo/gql-test-suite

# Set the working directory in the Docker image
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gts .

# Use a smaller image to run our application
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/gts .

# Expose port 8080 for HTTP traffic
EXPOSE 8080

ENTRYPOINT ["/app/gts"]

# Command to run the application
CMD ["crawl", "server", "start"]
