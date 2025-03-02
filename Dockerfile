FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/orchestrator ./cmd/orchestrator
RUN CGO_ENABLED=0 go build -o /out/agent        ./cmd/agent

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /out/orchestrator /app/orchestrator
COPY --from=builder /out/agent        /app/agent

COPY --from=builder /app/config /app/config

EXPOSE 8080

CMD ["/app/orchestrator"]