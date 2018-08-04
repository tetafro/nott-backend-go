FROM golang:1.10-alpine AS build

WORKDIR /go/src/github.com/tetafro/nott-backend-go

COPY . .

RUN go build -o ./bin/nott ./cmd/nott

FROM alpine:3.7

WORKDIR /app

COPY migrations migrations
COPY --from=build /go/src/github.com/tetafro/nott-backend-go/bin/nott /app/

RUN addgroup -S -g 5000 nott && \
    adduser -S -u 5000 -G nott nott && \
    chown -R nott:nott .

USER nott

EXPOSE 8080

CMD ["/app/nott"]
