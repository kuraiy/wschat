package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedis(rdb *redis.Client) *Redis {
	return &Redis{
		rdb: rdb,
		ctx: context.Background(),
	}
}

func (r *Redis) CheckToken(ref string) bool {
	exist, err := r.rdb.Exists(r.ctx, ref).Result()
	if err != nil {
		return false
	}

	if exist == 1 {
		return true
	}

	return false
}

func (r *Redis) InsertToken(ref string, ttl int) bool {
	if err := r.rdb.Set(r.ctx, ref, true, time.Duration(ttl)*time.Hour).Err(); err != nil {
		return false
	}
	return true
}

func (r *Redis) DeleteToken(ref string) {
	_ = r.rdb.Del(r.ctx, ref).Err()
}
