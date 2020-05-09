FROM golang:1.14-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /application ./cmd/server/main.go

FROM alpine

COPY --from=builder /application /app/application

ENV ENV=production

ENTRYPOINT ./app/application