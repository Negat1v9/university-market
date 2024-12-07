FROM golang:1.23.4-alpine3.20 as builder

WORKDIR /build

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -o /main cmd/marketplace/main.go

FROM alpine:latest

COPY --from=builder main /bin/main

ENTRYPOINT ["/bin/main"]