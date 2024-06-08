FROM golang:1.22.3-alpine AS builder

COPY . /github.com/passsquale/chat-server/source/
WORKDIR /github.com/passsquale/chat-server/source/

RUN go mod download
RUN go build -o ./bin/chat_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/passsquale/chat-server/source/bin/chat_server .

CMD ["./chat_server"]