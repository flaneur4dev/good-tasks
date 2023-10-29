# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/scheduler/ cmd/scheduler/
COPY internal/ internal/

RUN GOOS=linux go build -o /scheduler ./cmd/scheduler/

FROM alpine:3.18

COPY --from=build /scheduler /scheduler
COPY ./configs/scheduler.prod.yaml /etc/scheduler/scheduler.prod.yaml

CMD ["/scheduler", "-config", "/etc/scheduler/scheduler.prod.yaml"]
