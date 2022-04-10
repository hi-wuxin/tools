package rdb

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type rdbCmd interface {
	//redis.Cmdable
	baseCmd
}

type baseCmd interface {
	ZRangePop(ctx context.Context, key string, start, stop int64) (result []string, err error)
	ZRangeWithScoresPop(ctx context.Context, key string, start, stop int64) (result []redis.Z, err error)
	ZRangeByScorePop(ctx context.Context, key string, opt *redis.ZRangeBy) (result []string, err error)
	ZRangeByScoreWithScoresPop(ctx context.Context, key string, opt *redis.ZRangeBy) (result []redis.Z, err error)
}
