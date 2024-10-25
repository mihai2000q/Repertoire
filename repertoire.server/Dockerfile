############################
# STEP 1 build executable binary
############################
FROM golang:1.23 AS builder

# Set Working directory and Copy the application
WORKDIR /app
COPY . .

# Install dependencies
RUN go mod download

# Build the binary executable
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/main

############################
# STEP 2 build a small image
############################
FROM alpine:3

# Set Working directory
WORKDIR /

# Copy our static executable and environment file
COPY --from=builder /go/main /go/main
COPY --from=builder /app/.env /go/.env

# Expose default port
EXPOSE 8080

# Run
ENTRYPOINT ["/go/main"]