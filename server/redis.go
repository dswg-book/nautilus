package server

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Host, Port string
}

type Redis struct {
	Options *RedisOptions
	client  *redis.Client
	Context context.Context
}

var (
	redisInstance    *Redis
	defaultRedisHost = "127.0.0.1"
	defaultRedisPort = "6379"
)

func ensureClient(r *Redis) {
	if r == nil {
		r = NewRedis(context.Background(), RedisOptions{Host: defaultRedisHost, Port: defaultRedisPort})
	}
	if r.client == nil {
		client := redis.NewClient(&redis.Options{Addr: r.Options.Host + ":" + r.Options.Port})
		r.client = client
		return
	}
}

func NewRedis(ctx context.Context, options RedisOptions) *Redis {
	redisInstance = &Redis{
		Options: &options,
		Context: ctx,
	}
	ensureClient(redisInstance)
	return redisInstance
}

func RedisInstance() *Redis {
	return redisInstance
}

func (r *Redis) Connect() error {
	ensureClient(r)
	cmd := r.client.Ping(r.Context)
	return cmd.Err()
}

func (r *Redis) Close() error {
	ensureClient(r)
	return r.client.Close()
}
