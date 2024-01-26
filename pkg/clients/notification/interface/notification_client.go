package interfaces

import "github.com/ashiqsabith123/love-bytes-proto/notifications/pb"

type NotificationClient interface {
	InitNotificationClient()
	GetClient() pb.NotificationServiceClient
}
