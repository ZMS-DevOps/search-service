# Start from the latest golang base image
FROM golang:1.22.1 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./go.sum ./

# # Copy the local dependency
# COPY /common ../common

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the Working Directory inside the container
COPY .. .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .



######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
RUN chown root:root /root/main && chmod 755 /root/main
RUN chmod 555 /root/main
# Expose port 8000 to the outside world
EXPOSE 8083
USER appuser
# Command to run the executable
CMD ["./main"]