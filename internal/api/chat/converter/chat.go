package converter

import (
	"github.com/passsquale/chat-server/internal/model"
	chatPb "github.com/passsquale/chat-server/pkg/chat_v1"
)

func ProtoToMessage(message *chatPb.Message) model.MessageDTO {
	return model.MessageDTO{
		Author:    message.From,
		Content:   message.Text,
		CreatedAt: message.Timestamp.AsTime(),
	}
}
