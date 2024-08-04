package chatapi

import (
	"context"
	"fmt"

	"github.com/passsquale/chat-server/internal/model"
	chatPb "github.com/passsquale/chat-server/pkg/chat_v1"
)

func (i *Implementation) Create(ctx context.Context, req *chatPb.CreateRequest) (*chatPb.CreateResponse, error) {
	res, err := i.chatServ.Create(ctx, model.ChatDTO{
		Usernames: req.Usernames,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %v", err)
	}

	return &chatPb.CreateResponse{
		Id: res,
	}, nil
}
