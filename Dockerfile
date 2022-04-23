# Dockerfile References: https://docs.docker.com/engine/reference/builder/

FROM golang:1.15.11-buster as builder

LABEL maintainer="Huseyin Cakir <huseyin-cakir-2013@yandex.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags="-extldflags=-static" -o goauthz .


######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/goauthz .

EXPOSE 3000

CMD ["./goauthz"] 
