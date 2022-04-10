package rdb

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

var (
	rdbClient *Client
	ctx       context.Context
)

func init() {
	ctx = context.Background()
	rdbClient = NewClient(ctx, &Options{
		Network:  "tcp",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

}

func TestClient_ZRangeByScorePop(t *testing.T) {

}

func TestClient_ZRangeByScoreWithScoresPop(t *testing.T) {

}

func TestClient_ZRangePop(t *testing.T) {
	t.Log(rdbClient.ZRangePop(ctx, "zz", 0, 1))
}

func TestClient_ZRangeWithScoresPop(t *testing.T) {

}

func BenchmarkPop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rdbClient.EvalSha(ctx, zpopSha, []string{"zz"}, 0, 0).Result()
		//rdbClient.ZRange(ctx, "zz", 0, 10)
	}
}

func BenchmarkName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rdbClient.ZAdd(ctx, "zz", &redis.Z{Score: float64(i), Member: i})
	}
}
func TestName(t *testing.T) {

}
