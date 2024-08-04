package repository

import (
	"context"

	"github.com/passsquale/chat-server/internal/model"
)

type ChatRepository interface {
	Create(context.Context, model.ChatDTO) (int64, error)
	Delete(context.Context, int64) error
}

type MessagesRepository interface {
	Create(context.Context, model.MessageDTO) error
}

type LogsRepository interface {
	Create(ctx context.Context, log model.Log) (int64, error)
	Get(ctx context.Context, id int64) (model.Log, error)
}
