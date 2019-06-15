FROM golang:1.12-alpine AS build

WORKDIR /build

RUN apk add --no-cache git gcc musl-dev

COPY . .

RUN go mod download && \
    go build -o ./bin/nott ./cmd/nott

FROM alpine:3.9

WORKDIR /app

COPY migrations migrations
COPY --from=build /build/bin/nott /app/

RUN apk add --no-cache ca-certificates && \
    addgroup -S -g 5000 nott && \
    adduser -S -u 5000 -G nott nott && \
    chown -R nott:nott .

USER nott

EXPOSE 8080

CMD ["/app/nott"]
