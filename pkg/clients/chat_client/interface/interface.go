package interfaces

import "github.com/ashiqsabith123/love-bytes-proto/chat/pb"

type ChatClient interface {
	InitChatClient()
	GetClient() pb.ChatServiceClient
}
