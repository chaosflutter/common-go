# Build stage
FROM golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum ./

ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o main ./main.go

# Final stage
FROM alpine:latest
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Set environment variable
ENV ENV=production
RUN mkdir -p /root/logs

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
# Copy the public directory

EXPOSE 8080

VOLUME /root/logs

CMD ["./main"]
