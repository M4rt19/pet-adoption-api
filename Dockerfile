# ---- Build stage ----
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Enable automatic toolchain download so Go can satisfy go 1.25.1
ENV GOTOOLCHAIN=auto

# Install git (needed if modules pulled from git)
RUN apk add --no-cache git

# Copy go mod/sum first and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pet-adoption-api ./cmd/main.go

# ---- Run stage ----
FROM alpine:3.19

WORKDIR /app

RUN adduser -D appuser
USER appuser

COPY --from=builder /app/pet-adoption-api .

ENV SERVER_PORT=8080

EXPOSE 8080

CMD ["./pet-adoption-api"]
