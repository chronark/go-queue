package queue

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

type Queue struct {
	queueDone     string
	queueProgress string
	queueTodo     string
	redisClient   *redis.Client
	ctx           context.Context
}

func NewQueue(tenantId string, redisClient *redis.Client) (*Queue, error) {
	
	return &Queue{
		queueDone:     fmt.Sprintf("%s:done", tenantId),
		queueProgress: fmt.Sprintf("%s:inProgress", tenantId),
		queueTodo:     fmt.Sprintf("%s:todo", tenantId),
		redisClient:   redisClient,
		ctx:           context.Background(),
	}, nil
}

func (q *Queue) Produce(key string, message SignedMessage) error {
	serializedMessage, err := message.Serialize()
	if err != nil {
		return fmt.Errorf("Unable to serialize message: %w", err)
	}
	err = q.redisClient.LPush(q.ctx, fmt.Sprintf("%s:%s", q.queueTodo, key), serializedMessage).Err()
	if err != nil {
		return fmt.Errorf("Unable to push message to queue: %w", err)
	}
	return nil
}

func (q *Queue) Consume(key string) (string, error) {
	res := q.redisClient.RPopLPush(q.ctx, fmt.Sprintf("%s:%s", q.queueTodo, key), fmt.Sprintf("%s:%s", q.queueProgress, key))
	if res.Err() != nil {
		return "", res.Err()
	}
	return res.Val(), nil
}

func (q *Queue) Acknowledge(key string) error {
	res := q.redisClient.LMove(q.ctx, fmt.Sprintf("%s:%s", q.queueProgress, key), fmt.Sprintf("%s:%s", q.queueDone, key), "right", "right")
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
