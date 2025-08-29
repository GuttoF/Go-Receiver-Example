# -------------------------------------- Stage 1 --------------------------------------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server .

# -------------------------------------- Stage 2 --------------------------------------
FROM alpine:latest

WORKDIR /

COPY --from=build /app/server .
ENV PORT=8080
ENV FUNCTION_TARGET=ReceiverFunction

EXPOSE 8080

CMD ["./server"]