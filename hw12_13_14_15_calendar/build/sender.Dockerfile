# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/sender/ cmd/sender/
COPY internal/ internal/

RUN GOOS=linux go build -o /sender ./cmd/sender/

FROM alpine:3.18

COPY --from=build /sender /sender
COPY ./configs/sender.prod.yaml /etc/sender/sender.prod.yaml

CMD ["/sender", "-config", "/etc/sender/sender.prod.yaml"]
