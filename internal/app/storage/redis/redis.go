package RedisCache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wanna-beat-by-bit/goMetricService/internal/app/config"
)

type SessionCache struct {
	client *redis.Client
}

func New(conf *config.Config) (*SessionCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host,
		Password: conf.Redis.Password, // no password set
		DB:       conf.Redis.Name,     // use default DB
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("can't establish connection with redis: %s", err.Error())
	}

	return &SessionCache{
		client: client,
	}, nil
}

func (sc *SessionCache) GetSession(userId string) (string, error) {
	value, err := sc.client.Get(context.Background(), userId).Result() //Get session
	sessionId := value
	if err != redis.Nil && err != nil {
		return "", fmt.Errorf("error while getting key: %s", err.Error()) //unexpected error
	}
	if err == redis.Nil {
		sessionId = uuid.New().String()
		err = sc.client.SetEx(context.Background(), userId, sessionId, time.Second*60).Err() //Create new session
		if err != nil {
			return "", fmt.Errorf("error while creating key: %s", err.Error())
		}
		return sessionId, nil
	}

	err = sc.client.SetEx(context.Background(), userId, sessionId, time.Second*60).Err() //Update session if exists
	if err != nil {
		return "", fmt.Errorf("error while creating(updating) key: %s", err.Error())
	}

	return sessionId, nil
}
