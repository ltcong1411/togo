FROM golang:1.17-alpine3.10 as builder
WORKDIR  /home/togo
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ./togo


FROM alpine:3
RUN apk add --update ca-certificates
RUN apk add --no-cache tzdata && \
    cp -f /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime && \
    apk del tzdata
WORKDIR /app
COPY config/config.yaml config/
COPY --from=builder /home/togo .

ENTRYPOINT ["./togo"]
