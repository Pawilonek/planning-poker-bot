FROM golang:1.20-alpine AS builder

RUN mkdir /app
WORKDIR /app

ADD . /app

RUN go build -o bin/scrumpoke ./cmd/scrumpoke


FROM alpine:latest AS production
RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/bin /app

CMD ["./scrumpoke"]

