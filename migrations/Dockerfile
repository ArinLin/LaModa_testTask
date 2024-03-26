FROM golang:1.22-alpine

RUN go install github.com/jackc/tern/v2@latest

COPY sql /migrations

ENTRYPOINT ["tern"]