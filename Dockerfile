FROM golang:alpine AS builder
RUN apk add --no-cache upx
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -ldflags "-s -w" -a -installsuffix cgo -o app main.go
RUN upx --ultra-brute -qq app
RUN upx -t app

FROM scratch
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["/app"]
