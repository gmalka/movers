FROM golang:1.19 AS builder
WORKDIR /movers
COPY . .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./app ./cmd/main.go

FROM alpine:latest
WORKDIR /movers
COPY ./.env .
COPY ./templates ./templates
COPY --from=builder /movers/app .
ENTRYPOINT [ "/movers/app" ]