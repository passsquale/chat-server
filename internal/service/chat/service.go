package chatservice

import (
	"github.com/passsquale/chat-server/internal/repository"
	"github.com/passsquale/chat-server/internal/service"

	"github.com/passsquale/platform_common/pkg/client/db"
)

type serv struct {
	chatRepo  repository.ChatRepository
	txManager db.TxManager
	logsRepo  repository.LogsRepository
}

func NewService(chatRepo repository.ChatRepository, tx db.TxManager, logsRepo repository.LogsRepository) service.ChatService {
	return &serv{
		chatRepo:  chatRepo,
		txManager: tx,
		logsRepo:  logsRepo,
	}
}
