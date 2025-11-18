ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/todo-list-be ./cmd/server

FROM alpine:3.19
WORKDIR /app

RUN addgroup -S app && adduser -S app -G app && \
	apk add --no-cache ca-certificates

COPY --from=builder /bin/todo-list-be /app/todo-list-be

EXPOSE 8080
USER app
CMD ["/app/todo-list-be"]
