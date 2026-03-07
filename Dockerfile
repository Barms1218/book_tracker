FROM golang:1.25-alpine AS builder

# Install build dependencies for CGO
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY . .

# Enable CGO and build
RUN CGO_ENABLED=1 GOOS=linux go build -o /book_tracker

FROM alpine:latest

ENV DATABASE_URL=/data/book_database.db

RUN apk add --no-cache ca-certificates sqlite-libs


WORKDIR /root/
COPY --from=builder /book_tracker .
COPY index.html .

RUN mkdir /data

VOLUME /data
EXPOSE 8080

CMD ["./book_tracker"]

