package cache

import (
	"WalletTopUp/internal/config"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func Connect(ctx context.Context, conf config.CacheClient) *redis.Client {
	addr := fmt.Sprintf("%s:%s", conf.Host, conf.Port)
	client := redis.NewClient(&redis.Options{Addr: addr, Password: conf.Pass, DB: conf.DB})
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		log.Println("failed to connect redis.")
		return nil
	}

	log.Println("Redis connected.")
	return client
}

func Close(client *redis.Client) error {
	if client == nil {
		return errors.New("MemDB didn't initialize.")
	}

	return client.Close()
}
