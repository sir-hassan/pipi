FROM golang:1.16-buster

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go get -u golang.org/x/lint/golint
COPY . .
RUN make build

WORKDIR /app
ENTRYPOINT ["/app/bin/pipi"]