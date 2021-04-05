FROM golang:1.16-buster

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build

COPY --chown=0:0 --from=builder /app/bin/pipi /app/pipi

WORKDIR /app
ENTRYPOINT ["/app/pipi"]