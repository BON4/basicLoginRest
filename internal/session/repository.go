package session

import (
	"context"
	"time"
)

type Repository interface {
	//Check if cache exists
	Check(ctx context.Context, sid string) (bool, error)

	//Create create cache and specify expiration time
	Create(ctx context.Context, sid string, expired time.Duration) (Store, error)

	//Update gets value from cache storage, update its ttl, and load it in local storage
	Update(ctx context.Context, sid string, expired time.Duration) (Store, error)

	//Delete deletes val from cache storage
	Delete(ctx context.Context, sid string) error

	//Refresh Use sid to replace old sid and return session store
	Refresh(ctx context.Context, oldSid, sid string, expired time.Duration) (Store, error)
}

type Store interface {
	Context() context.Context
	SessionID() string
	//Set set cache
	Set(key string, value interface{})
	//Get get cache
	Get(key string) (interface{}, bool)
	//Delete deletes key from store
	Delete(key string) interface{}
	//Save saves cache in to store
	Save() error
	//Flush flushes cache in store
	Flush() error
}
