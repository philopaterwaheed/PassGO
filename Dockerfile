FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /app/passgo-backend ./cmd/backend

FROM alpine:3.22

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/passgo-backend ./passgo-backend

ENV GIN_MODE=release
ENV PORT=8080

EXPOSE 8080

CMD ["./passgo-backend"]
