package api

import "github.com/passsquale/chat-server/internal/model"

var RouteAccesses = map[string][]model.UserRole{
	"/chat_v1.ChatV1/SendMessage": {model.USER},
	"/chat_v1.ChatV1/Create":      {model.USER},
	"/chat_v1.ChatV1/Delete":      {model.ADMIN},
}
