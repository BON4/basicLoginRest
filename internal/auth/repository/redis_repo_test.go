package repository

import (
	"basicLoginRest/internal/auth"
	"basicLoginRest/internal/models"
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func SetUpRedis() auth.KVRepository {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return NewRedisRepository(client)
}

func TestAuthRedisRepo_SetUser(t *testing.T) {
	t.Parallel()

	repo := SetUpRedis()

	t.Run("SetUser", func(t *testing.T) {
		key := 0
		usr, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN, "1234")
		err := repo.SetUser(context.Background(), strconv.Itoa(key), 1, &usr)
		require.NoError(t, err)
	})
}

func TestAuthRedisRepo_DeleteUser(t *testing.T) {
	t.Parallel()
	repo := SetUpRedis()

	t.Run("DeleteUser", func(t *testing.T) {
		key := 0
		err := repo.DeleteUser(context.Background(), strconv.Itoa(key))
		require.NoError(t, err)
	})
}

func TestAuthRedisRepo_GetUserByID(t *testing.T) {
	t.Parallel()
	repo := SetUpRedis()

	t.Run("DeleteUser", func(t *testing.T) {
		key := strconv.Itoa(0)
		usr, _ := userFactory.NewUser("abcd", "abcd@email.com", models.ADMIN, "1234")
		err := repo.SetUser(context.Background(), key, 1, &usr)
		require.NoError(t, err)

		getUser, err := repo.GetUserByID(context.Background(), key)
		require.NoError(t, err)
		require.Equal(t, usr, *getUser)
	})
}