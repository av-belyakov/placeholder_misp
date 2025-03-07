package redisapi_test

import (
	"context"
	"fmt"
	"testing"

	redis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const (
	HostRDb string = "127.0.0.1"
	PortRDb int    = 6379
)

func TestGetDumpdb(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", HostRDb, PortRDb),
	})

	status := client.Ping(context.Background())
	res, err := status.Result()
	assert.NoError(t, err)

	assert.NotEmpty(t, res)
	assert.Equal(t, res, "PONG")

	//list := client.HGetAll(context.Background(), "")

}
