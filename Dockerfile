FROM golang:1.23-alpine AS build

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with CGO enabled for SQLite
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server main.go

FROM alpine:latest

WORKDIR /app

# Install runtime dependencies for SQLite
RUN apk add --no-cache ca-certificates sqlite-libs

COPY --from=build /app/server /app/server
COPY --from=build /app/web /app/web

# Ensure the database file exists or is in a volume
# EXPOSE 30022 (Match your main.go port if needed)
EXPOSE 30022

CMD ["/app/server"]
