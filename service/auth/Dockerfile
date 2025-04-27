    FROM golang:1.23-alpine AS build  

    WORKDIR /app
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    RUN CGO_ENABLED=0 go build -o /auth ./cmd
    
    FROM alpine:latest                
    COPY --from=build /auth /usr/local/bin/auth 
    ENTRYPOINT ["/usr/local/bin/auth"]