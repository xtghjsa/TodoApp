package pkg

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitializeRedis(redisHost, redisPort string) (rds *redis.Client, ctx context.Context, err error) {
	ctx = context.Background()
	rds = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})
	_, err = rds.Ping(ctx).Result()
	if err != nil {
		return nil, nil, err
	}
	return rds, ctx, nil
}
