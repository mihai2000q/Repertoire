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

# Expose port
ARG SERVER_PORT
EXPOSE $SERVER_PORT

# Run
ENTRYPOINT ["/go/main"]