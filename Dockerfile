# ----- BUILD STAGE -----

FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o main

# ----- RUN STAGE -----

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/static ./static
COPY --from=builder /app/.env .env
COPY --from=builder /app/main .

RUN chmod +x main

EXPOSE 8080

CMD ["./main"]
