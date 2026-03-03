FROM golang:1.25-alpine AS builder

# Install build dependencies for CGO
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY . .

# Enable CGO and build
RUN CGO_ENABLED=1 GOOS=linux go build -o /book_tracker

FROM alpine:latest

RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /root/
COPY --from=builder /book_tracker .

VOLUME /data

CMD ["./book_tracker"]

