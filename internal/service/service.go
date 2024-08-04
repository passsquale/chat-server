package service

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -o ./mocks/ -s ".go"

import (
	"context"

	"github.com/passsquale/chat-server/internal/model"
)

type ChatService interface {
	Create(context.Context, model.ChatDTO) (int64, error)
	Delete(context.Context, int64) error
}

type MessageService interface {
	SendMessage(context.Context, model.MessageDTO) error
}
