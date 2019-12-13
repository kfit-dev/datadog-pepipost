FROM golang:latest AS builder

WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates dumb-init
RUN addgroup -S app
RUN adduser -S app -G app
USER app
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 8080
ENTRYPOINT ["/usr/bin/dumb-init","--"]
CMD ["./app"]
