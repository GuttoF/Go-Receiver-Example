# -------------------------------------- Stage 1 --------------------------------------
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# -------------------------------------- Stage 2 --------------------------------------
FROM alpine:latest

WORKDIR /

COPY --from=builder /app/server .
RUN chmod +x server

ENV PORT=8080

EXPOSE 8080

CMD ["./server"]