# Use a Go image to build the application
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# CGO_ENABLED=0 is important for static linking, making the binary self-contained
# -o specifies the output file name
# ./cmd/controller is the path to the main package
RUN mkdir -p bin && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/controller ./cmd/controller

# Use a minimal base image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app/bin/controller /bin/controller

# Command to run the application
ENTRYPOINT ["/bin/controller"]
