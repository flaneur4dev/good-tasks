# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/calendar/ cmd/calendar/
COPY internal/ internal/

RUN GOOS=linux go build -o /calendar ./cmd/calendar/

FROM alpine:3.18

COPY --from=build /calendar /calendar
COPY ./configs/calendar.prod.yaml /etc/calendar/calendar.prod.yaml

EXPOSE 3000
EXPOSE 50051

CMD ["/calendar", "-config", "/etc/calendar/calendar.prod.yaml"]
