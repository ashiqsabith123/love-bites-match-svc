package chat

import (
	logs "github.com/ashiqsabith123/love-bytes-proto/log"

	"github.com/ashiqsabith123/love-bytes-proto/chat/pb"
	interfaces "github.com/ashiqsabith123/match-svc/pkg/clients/chat_client/interface"
	"github.com/ashiqsabith123/match-svc/pkg/config"
	"google.golang.org/grpc"
)

type ChatClient struct {
	config config.Config
}

var Conn *grpc.ClientConn
var err error

func NewChatClient(config config.Config) interfaces.ChatClient {
	client := &ChatClient{config: config}
	client.InitChatClient()
	return &ChatClient{}
}

func (A *ChatClient) InitChatClient() {

	// credentials, err := helper.GetCertificate("pkg/services/auth-svc/cert/ca-cert.pem", "pkg/services/auth-svc/cert/client-cert.pem", "pkg/services/auth-svc/cert/client-key.pem")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	Conn, err = grpc.Dial(A.config.Port.ChatSvcPort, grpc.WithInsecure())
	if err != nil {
		logs.ErrLog.Println("Could not connect the chat server:", err)
	}

	logs.GenLog.Println("Chat service connected at port ", A.config.Port.ChatSvcPort)

}

func (A *ChatClient) GetClient() pb.ChatServiceClient {
	return pb.NewChatServiceClient(Conn)
}
