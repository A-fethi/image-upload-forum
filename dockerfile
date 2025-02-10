FROM golang:1.23-alpine AS builder

LABEL \
    website.name="Forum" \
    description="A web forum application built with Go and SQLite and JavaScript. It allows user communication through posts and comments, supports authentication with sessions and cookies, implements a like/dislike system, and provides filtering options. The application follows best practices, handles HTTP and technical errors, and is containerized with Docker for easy deployment." \
    authors="Fethi Abderrahmane, Aymane Berhili, Mostafa Zakri, Anass Elabsi, Jamal Bajady"

WORKDIR /app

RUN apk add --no-cache gcc g++ musl-dev sqlite-dev

ENV CGO_ENABLED=1

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./app/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/app/db/schema.sql ./app/db/schema.sql
COPY --from=builder /app/logs/ ./logs/
COPY --from=builder /app/static/ ./static/
COPY --from=builder /app/templates/ ./templates/

RUN apk add --no-cache sqlite-dev bash

EXPOSE 8080

CMD ["./main"]
