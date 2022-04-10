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
	ZRangePop(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	ZRangeWithScoresPop(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd
	ZRangeByScorePop(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd
}
