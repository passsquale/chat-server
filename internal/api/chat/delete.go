package chatapi

import (
	"context"
	"fmt"

	chatPb "github.com/passsquale/chat-server/pkg/chat_v1"

	"github.com/golang/protobuf/ptypes/empty"
)

func (i *Implementation) Delete(ctx context.Context, req *chatPb.DeleteRequest) (*empty.Empty, error) {
	err := i.chatServ.Delete(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete chat: %v", err)
	}

	return &empty.Empty{}, nil
}
