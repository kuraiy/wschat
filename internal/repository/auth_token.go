package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
}

func NewRedis(rdb *redis.Client) *Redis {
	return &Redis{
		rdb: rdb,
	}
}

func (r *Redis) CheckToken(ctx context.Context, ref string) bool {
	exist, err := r.rdb.Exists(ctx, ref).Result()
	if err != nil {
		return false
	}

	if exist == 1 {
		return true
	}

	return false
}

func (r *Redis) GetToken(ctx context.Context, ref string) (int64, error) {
	idStr, err := r.rdb.Get(ctx, ref).Result()

	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Redis) InsertToken(ctx context.Context, ref string, id int64, ttl int) bool {
	if err := r.rdb.Set(ctx, ref, id, time.Duration(ttl)*time.Hour).Err(); err != nil {
		return false
	}
	return true
}

func (r *Redis) DeleteToken(ctx context.Context, ref string) {
	_ = r.rdb.Del(ctx, ref).Err()
}
