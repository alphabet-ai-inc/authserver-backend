# Example: https://docs.docker.com/guides/golang/build-images/
FROM golang:1.24-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o goapps .

# Stage 2: Create minimal runtime image
FROM alpine:3.22.0
WORKDIR /app
COPY --from=build /app/goapps .
COPY .env .

EXPOSE 8080
CMD ["./goapps"]
# CMD ["/bin/sh", "-c", "sleep infinity"]