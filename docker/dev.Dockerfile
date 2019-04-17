FROM golang:alpine

ENV GO111MODULE on
ARG PROJECT=godev

WORKDIR /go/src/${PROJECT}/

RUN apk add ca-certificates ffmpeg libva-intel-driver git gcc musl-dev && \
    apk upgrade

COPY go.mod go.sum ./
RUN go mod download && \
    go get -u github.com/enfipy/gouto

CMD ["gouto", "-dir=src"]
