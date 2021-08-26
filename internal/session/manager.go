package session

import "context"

type Manager interface {
	Start(ctx context.Context, sid string) (Session, error)
	Refresh(ctx context.Context, oldSid string) (Session, error)
}