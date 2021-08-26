package repository

import (
	"basicLoginRest/internal/auth"
	"basicLoginRest/internal/models"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"time"
)

type authRedisRepo struct {
	redisClient *redis.Client
}

func (r *authRedisRepo) GetUserByKey(ctx context.Context, key string) (*models.User, error) {
	binaryData, err := r.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "authRedisRepo.GetUserByID.GetBytes")
	}

	var user models.User
	if err = json.Unmarshal(binaryData, &user); err != nil {
		return nil, errors.Wrap(err, "authRedisRepo.GetUserByID.Unmarshal")
	}
	return &user, nil
}

func (r *authRedisRepo) SetUser(ctx context.Context, key string, t int, user *models.User) error {
	marshUser, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "authRedisRepo.SetUser.Marshal")
	}

	if err = r.redisClient.SetEX(ctx, key, marshUser, time.Duration(t)*time.Second).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo.SetUser.SetEX")
	}
	return nil
}

func (r *authRedisRepo) DeleteUser(ctx context.Context, key string) error {
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		return errors.Wrap(err, "authRedisRepo.DeleteUser.Del")
	}
	return nil
}

func NewRedisRepository(rc *redis.Client) auth.KVRepository {
	return &authRedisRepo{
		redisClient: rc,
	}
}
