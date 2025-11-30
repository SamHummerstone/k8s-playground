# Multi-stage Dockerfile for a Go cookiecutter app
ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /cookie-cutter

# required to fetch some modules and provide CA certs for the final image
RUN apk add --no-cache git ca-certificates

# download modules (speed up rebuilds)
COPY go.mod go.sum ./
RUN go mod download

# build the binary (assumes main package is at project root; adjust path if needed)
COPY . .
ARG CGO_ENABLED=0
ENV CGO_ENABLED=${CGO_ENABLED} GOOS=linux GOARCH=arm64
RUN go build -ldflags="-s -w" -o /app .

# Final minimal image
FROM gcr.io/distroless/static:nonroot
COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER nonroot
EXPOSE 8080
ENTRYPOINT ["/app"]