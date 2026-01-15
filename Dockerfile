
# Build stage
FROM golang:1.24-alpine AS builder

# Install CA certs & git
RUN apk add --no-cache ca-certificates git

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o portfolio-service .


# Runtime stage using distroless image
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/portfolio-service /app/portfolio-service

EXPOSE 8080

ENTRYPOINT ["/app/portfolio-service"]
CMD ["server"]
