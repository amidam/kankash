package webhook

import (
	"context"

	"kankash/internal/event/push"
)

type service struct {
	
}

type Service interface {
	Push(ctx context.Context, req push.Payload) error
}

