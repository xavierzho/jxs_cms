package redisdb

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(host, password string, db int) (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		rdb,
	}, nil
}

// Lock TODO 锁通过 第三方包来实现
func (r *RedisClient) Lock(ctx context.Context, key string) (err error) {
	lockRKey := getLockRKey(key)
	unlockChannel := getUnlockRKey(key)
	getLockFunc := func() (bool, error) {
		getLock, err := r.SetNX(ctx, lockRKey, 1, time.Second*30).Result()
		if err != nil {
			return false, err
		}
		if getLock {
			// 解锁
			return true, nil
		} else {
			pubSub := r.Subscribe(ctx, unlockChannel)
			defer pubSub.Close()
			// 十分钟后若还未获取到锁则报错
			timer := time.NewTicker(time.Minute * 3)
			defer timer.Stop()

			select {
			case <-pubSub.Channel():
				// 收到解锁消息后获取锁
				getLock, err = r.SetNX(ctx, lockRKey, 1, time.Second*30).Result()
				if err != nil {
					return false, err
				}
				if getLock {
					// 解锁
					return true, nil
				}
			case <-timer.C:
				return false, fmt.Errorf("redis get lock fail[%s]", lockRKey)
			}
		}
		return false, err
	}

	var isGetLock bool
	for {
		isGetLock, err = getLockFunc()
		if err != nil {
			return err
		}
		if isGetLock {
			break
		}
	}
	return err
}

func (r *RedisClient) Unlock(ctx context.Context, key string) (err error) {
	lockRKey := getLockRKey(key)
	unlockChannel := getUnlockRKey(key)
	_, err = r.Del(ctx, lockRKey).Result()
	if err != nil {
		return err
	}
	// 发送解锁消息
	_, err = r.Publish(ctx, unlockChannel, 1).Result()
	if err != nil {
		return err
	}
	return nil
}

func getLockRKey(key string) string {
	return fmt.Sprintf("Lock:%s", key)
}

func getUnlockRKey(key string) string {
	return fmt.Sprintf("Unlock:%s", key)
}
