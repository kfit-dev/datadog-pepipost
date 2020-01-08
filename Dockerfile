# build app
FROM golang:alpine AS builder
RUN apk add --no-cache upx
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOPROXY=https://proxy.golang.org go build -ldflags "-s -w" -a -installsuffix cgo -o app main.go
RUN upx --ultra-brute -qq app
RUN upx -t app

# build tzdata
FROM alpine:latest as alpine
RUN apk --no-cache add tzdata zip ca-certificates
WORKDIR /usr/share/zoneinfo
# -0 means no compression.  Needed because go's
# tz loader doesn't handle compressed data.
RUN zip -r -0 /zoneinfo.zip .


FROM scratch
ENV TZ=Asia/Kuala_Lumpur
ENV ZONEINFO /zoneinfo.zip
COPY --from=alpine /zoneinfo.zip /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["/app"]
