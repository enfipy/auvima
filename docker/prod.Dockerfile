# Build stage
FROM golang:1.11-alpine AS builder

ENV GO111MODULE on
ARG PROJECT=godev

WORKDIR /go/src/${PROJECT}/

COPY go.mod go.sum ./
RUN apk add ca-certificates ffmpeg libva-intel-driver git gcc musl-dev && \
    apk upgrade && \
    go mod download
COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app ./src

# Final stage
FROM jrottenberg/ffmpeg:scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app .
ENTRYPOINT ["./app"]
