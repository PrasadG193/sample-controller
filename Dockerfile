# Development image
FROM golang:1.13-alpine3.10 AS BUILD-ENV

ARG GOOS_VAL 
ARG GOARCH_VAL

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN GOOS=${GOOS_VAL} GOARCH=${GOARCH_VAL} go build -o /go/bin/controller .

# Production image
FROM alpine:3.10

# Create Non Privilaged user
COPY --from=BUILD-ENV /go/bin/controller /go/bin/controller

ENTRYPOINT /go/bin/controller
