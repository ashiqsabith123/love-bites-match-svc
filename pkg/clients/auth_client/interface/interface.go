package interfaces

import "github.com/ashiqsabith123/love-bytes-proto/auth/pb"

type AuthClient interface {
	InitAuthClient()
	GetClient() pb.AuthServiceClient
}
