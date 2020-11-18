package cacheClient

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type redisClient struct {
	*redis.Client
}

func newRedisClient(addr string) Client {
	c := redis.NewClient(&redis.Options{
		Addr: addr + ":7379",
	})
	_, err := c.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return &redisClient{c}
}

func (this *redisClient) get(key string) (string, error) {
	res, err := this.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return res, err
}

func (this *redisClient) set(key, value string) error {
	return this.Set(context.Background(), key, value, 0).Err()
}

func (this *redisClient) del(key string) error {
	return this.Del(context.Background(), key).Err()
}

func (this *redisClient) Run(cmd *Cmd) {

	if cmd.OpName == "get" {
		cmd.Value, cmd.Error = this.get(cmd.Key)
		return
	}

	if cmd.OpName == "set" {
		cmd.Error = this.set(cmd.Key, cmd.Value)
		return
	}

	if cmd.OpName == "del" {
		cmd.Error = this.del(cmd.Key)
		return
	}

	panic("unknown cmd name " + cmd.OpName)
}

func (this *redisClient) PipelineRun(cmds []*Cmd) {
	if len(cmds) == 0 {
		return
	}
	pipe := this.Pipeline()
	cmder := make([]redis.Cmder, len(cmds))
	for i, c := range cmds {
		if c.OpName == "set" {
			cmder[i] = pipe.Set(context.Background(), c.Key, c.Value, 0)
		} else if c.OpName == "get" {
			cmder[i] = pipe.Get(context.Background(), c.Key)
		} else if c.OpName == "del" {
			cmder[i] = pipe.Del(context.Background(), c.Key)
		} else {
			panic("unknown cmd name " + c.OpName)
		}
	}
	_, err := pipe.Exec(context.Background())
	if err != nil && err != redis.Nil {
		panic(err)
	}

	for i, c := range cmds {
		if c.OpName == "get" {
			value, e := cmder[i].(*redis.StringCmd).Result()
			if e == redis.Nil {
				value, e = "", nil
			}
			c.Value, c.Error = value, e
		} else {
			c.Error = cmder[i].Err()
		}
	}
}
