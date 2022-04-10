package rdb

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

type Options redis.Options
type Client struct {
	rdbCmd
	*redis.Client
}

func NewClient(ctx context.Context, opts *Options) *Client {
	client := redis.NewClient((*redis.Options)(opts))
	if err := client.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	zpopScoreSha, _ = client.ScriptLoad(ctx, popscore).Result()
	zpopSha, _ = client.ScriptLoad(ctx, pop).Result()
	zByScoreWithSha, _ = client.ScriptLoad(ctx, zrangeByScoreWith).Result()
	zByScoreSha, _ = client.ScriptLoad(ctx, zrangeByScore).Result()
	return &Client{
		Client: client,
	}
}

func (client *Client) ZRangePop(ctx context.Context, key string, start, stop int64) (result []string, err error) {

	return client.EvalSha(ctx, zpopSha, []string{key}, start, stop).StringSlice()

}

func (client *Client) ZRangeWithScoresPop(ctx context.Context, key string, start, stop int64) (result []redis.Z, err error) {

	data, err := client.EvalSha(ctx, zpopScoreSha, []string{key}, start, stop).StringSlice()
	if err != nil {
		return
	}
	return client.withScore(data)
}

func (client *Client) ZRangeByScorePop(ctx context.Context, key string, opt *redis.ZRangeBy) (result []string, err error) {
	return client.EvalSha(ctx, zByScoreSha, []string{key}, opt.Min, opt.Max, opt.Offset, opt.Count).StringSlice()
}
func (client *Client) ZRangeByScoreWithScoresPop(ctx context.Context, key string, opt *redis.ZRangeBy) (result []redis.Z, err error) {
	data, err := client.EvalSha(ctx, zByScoreWithSha, []string{key}, opt.Min, opt.Max, opt.Offset, opt.Count).StringSlice()
	if err != nil {
		return
	}
	return client.withScore(data)
}

func (client *Client) withScore(data []string) (result []redis.Z, err error) {
	if len(data) < 0 {
		return
	}
	if len(data)%2 != 0 {
		err = errors.New("data not aligned")
		return
	}
	result = make([]redis.Z, 0, len(data)/2)
	for i := 0; i < len(data); i++ {
		element := redis.Z{
			Member: data[i],
		}
		i++
		element.Score, _ = strconv.ParseFloat(data[i], 64)
		result = append(result, element)
	}
	return
}
