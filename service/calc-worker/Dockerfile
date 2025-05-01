FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/calc-worker ./cmd

FROM alpine:3.19
COPY --from=build /bin/calc-worker /usr/local/bin
ENTRYPOINT ["calc-worker"]