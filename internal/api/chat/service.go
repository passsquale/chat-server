package chatapi

import (
	"github.com/passsquale/chat-server/internal/service"
	chatPb "github.com/passsquale/chat-server/pkg/chat_v1"
)

type Implementation struct {
	chatPb.UnimplementedChatV1Server
	chatServ    service.ChatService
	messageServ service.MessageService
}

func NewImplementation(chatServ service.ChatService, messageServ service.MessageService) *Implementation {
	return &Implementation{
		chatServ:    chatServ,
		messageServ: messageServ,
	}
}
