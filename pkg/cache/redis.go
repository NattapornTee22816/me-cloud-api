package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

const (
	DBSession = 0
	DBCache   = 1
)

func NewCache(target int) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         os.Getenv("REDIS_ADDR"),
		Username:     os.Getenv("REDIS_USERNAME"),
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           target,
		PoolFIFO:     true,
		MaxIdleConns: 5,
		MaxRetries:   3,
		//TLSConfig: &tls.Config{
		//	InsecureSkipVerify: true,
		//},
	})

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctxTimeout).Err(); err != nil {
		return nil, err
	}

	return client, nil
}
