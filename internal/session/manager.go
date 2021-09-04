package session

import "context"

type Manager interface {
	Start(ctx context.Context, sid string) (Store, error)
	Refresh(ctx context.Context, oldSid string) (Store, error)
	Destroy(ctx context.Context, sid string) error
}