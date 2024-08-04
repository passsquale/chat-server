FROM golang:1.20.3-alpine as builder

COPY . /github.com/passsquale/chat-server/source
WORKDIR /github.com/passsquale/chat-server/source/

RUN go mod download
RUN go build -o ./bin/chat_server cmd/server/main.go

FROM alpine:latest

COPY --from=builder /github.com/passsquale/chat-server/source/bin/chat_server .
COPY --from=builder /github.com/passsquale/chat-server/source/prod.env .

CMD [ "./chat_server" ]