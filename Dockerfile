# syntax=docker/dockerfile:1.4
FROM golang:1.24 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /pikboard-api

FROM alpine:3.18

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /pikboard-api /pikboard-api

USER appuser

EXPOSE 8080

ENTRYPOINT ["/pikboard-api"]
