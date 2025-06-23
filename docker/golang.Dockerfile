FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../go.mod go.sum ./

RUN go mod download

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

FROM alpine:3.22.0

WORKDIR /app

RUN apk --no-cache add ca-certificates=20241121-r2 tzdata=2025b-r0

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
