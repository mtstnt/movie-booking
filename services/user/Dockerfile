# Builder
FROM golang:1.19.3-alpine3.16 as builder
WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download -x
COPY . .
RUN go build -o user .

# Runner
FROM alpine:3.16
WORKDIR /
COPY --from=builder /app/user ./user
EXPOSE 8000
ENTRYPOINT ["./user"]