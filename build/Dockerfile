FROM golang:latest AS builder
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gophermart ./cmd/gophermart/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /gophermart ./
RUN chmod +x ./gophermart
ENTRYPOINT ["./gophermart"]
EXPOSE 8080