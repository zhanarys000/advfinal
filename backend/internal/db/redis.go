package db

import (
	"adv/internal/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func SaveVerificationCodeToRedis(ctx context.Context, config config.Config, email, verificationCode string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.Redispassword,
		DB:       config.RedisDb,
	})
	defer rdb.Close()

	err := rdb.Set(ctx, fmt.Sprintf("verification_code:%s", email), verificationCode, time.Minute*10).Err()
	if err != nil {
		return err
	}
	return nil
}

func CheckVerificationCode(ctx context.Context, config config.Config, email, verificationCode string) (bool, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.Redispassword,
		DB:       config.RedisDb,
	})
	defer rdb.Close()

	val, err := rdb.Get(ctx, fmt.Sprintf("verification_code:%s", email)).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	if val == verificationCode {
		_ = rdb.Del(ctx, fmt.Sprintf("verification_code:%s", email))
		return true, nil
	}

	return false, nil
}
