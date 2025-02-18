package service

import "context"

type Realtime interface {
	PublishMessage(ctx context.Context, message string) error
	SubscribeToMessages(ctx context.Context, handler func(message string) error) error
}
