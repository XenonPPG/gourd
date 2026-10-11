# Stage 1: Build
FROM golang:1.27-alpine AS builder

COPY config.yaml /config.yaml

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/gourd/main.go

# Stage 2: Final image
FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /main .
COPY --from=builder /config.yaml .

EXPOSE 3000

CMD ["./main"]