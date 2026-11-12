FROM golang:1.27.1-alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go build -o bookmark_service cmd/api/main.go

FROM alpine:3.23.2 AS run

WORKDIR /opt/app

COPY --from=base /opt/app/bookmark_service /opt/app/bookmark_service

CMD ["/opt/app/bookmark_service"]