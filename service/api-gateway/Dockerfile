FROM golang:1.23 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/api-gateway ./cmd

FROM alpine:3.19
COPY --from=build /bin/api-gateway /usr/local/bin
ENTRYPOINT ["api-gateway"]