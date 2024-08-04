package messageservice

import (
	"github.com/passsquale/chat-server/internal/repository"
	"github.com/passsquale/chat-server/internal/service"

	"github.com/passsquale/platform_common/pkg/client/db"
)

type serv struct {
	messageRepo repository.MessagesRepository
	txManager   db.TxManager
	logsRepo    repository.LogsRepository
}

func NewService(messageRepo repository.MessagesRepository, tx db.TxManager, logsRepo repository.LogsRepository) service.MessageService {
	return &serv{
		messageRepo: messageRepo,
		txManager:   tx,
		logsRepo:    logsRepo,
	}
}
