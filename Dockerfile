# Step 1: Use the official Golang image as the base image
FROM golang:1.21 AS builder

# Step 2: Set the Current Working Directory inside the container
WORKDIR /app

# Step 3: Copy the Go Modules manifests
# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./

# Step 4: Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Step 5: Copy the entire project into the container's working directory
COPY . .

# Step 6: Build the Go app
RUN go build -o webhook-server .

# Step 7: Use a minimal base image for the final image
FROM alpine:latest

# Step 8: Install necessary dependencies (for example, certificates and libc)
#RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
RUN apk update && apk add --no-cache ca-certificates && \
    apk add --no-cache libc6-compat
# Step 9: Set the Current Working Directory inside the container
WORKDIR /app/

# Step 10: Copy the Go binary and the config file from the builder stage
COPY --from=builder /app/webhook-server /app/webhook-server
COPY config.json /app/config.json
RUN ls -la /app
# Step 11: Expose the port the app will run on
EXPOSE 8080

# Step 12: Command to run the application
CMD ["/app/webhook-server"]
