FROM golang:alpine AS builder
WORKDIR /app/
COPY . ./
RUN go build -o server -v ./cmd/server

FROM alpine
WORKDIR /app/
COPY --from=builder /app/server ./
COPY --from=builder /app/configs ./configs
COPY --from=builder /app/Makefile ./
CMD ["./server configs/dockerserver.yaml"]
