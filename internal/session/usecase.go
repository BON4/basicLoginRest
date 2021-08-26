package session

import (
	"context"
	"time"
)

type UCSession interface {
	//Check if cache exists
	Check(ctx context.Context, sid string) (bool, error)

	//Create create cache and specify expiration time
	Create(ctx context.Context, sid string, expired time.Duration) (Session, error)

	//Update gets value from cache storage, update its ttl, and load it in local storage
	Update(ctx context.Context, sid string, expired time.Duration) (Session, error)

	//Refresh Use sid to replace old sid and return session store
	Refresh(ctx context.Context, oldSid, sid string, expired time.Duration) (Session, error)
}

type Session interface {
	Context() context.Context

	SessionID() string
	//Set set cache
	Set(key string, value interface{})
	//Get get cache
	Get(key string) (interface{}, bool)
	//Save saves cache in to store
	Save() error
	//Flush flushes cache in store
	Flush() error
}