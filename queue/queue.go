package queue

import (
	"context"
	"fmt"

	redis "github.com/go-redis/redis/v8"
)

const (
	TODO        = "todo"
	IN_PROGRESS = "progress"
	DONE        = "done"
)

type Queue struct {
	redisClient *redis.Client
	ctx         context.Context
	tenantId    string
}

func NewQueue(tenantId string, redisClient *redis.Client) (*Queue, error) {

	return &Queue{
		redisClient: redisClient,
		ctx:         context.Background(),
		tenantId:    tenantId,
	}, nil
}

func (q *Queue) buildListKey(topic, status string) string {
	return fmt.Sprintf("%s:%s:%s", q.tenantId, topic, status)
}
func (q *Queue) buildKey(topic, status, id string) string {
	return fmt.Sprintf("%s:%s:%s:%s", q.tenantId, topic, status, id)
}

func (q *Queue) Produce(topic string, message SignedMessage) error {
	serializedMessage, err := message.Serialize()
	if err != nil {
		return fmt.Errorf("Unable to serialize message: %w", err)
	}

	err = q.redisClient.Set(q.ctx, fmt.Sprintf("%s:active:%s", q.tenantId, message.Message.Header.Id), serializedMessage, -1).Err()
	if err != nil {
		return err
	}

	err = q.redisClient.LPush(q.ctx, q.buildListKey(topic, TODO), message.Message.Header.Id).Err()
	if err != nil {
		return fmt.Errorf("Unable to push message to queue: %w", err)
	}
	return nil
}

func (q *Queue) Consume(topic string) (string, error) {
	messageId, err := q.redisClient.RPopLPush(q.ctx, q.buildListKey(topic, TODO), q.buildListKey(topic, IN_PROGRESS)).Result()
	if err != nil {
		return "", err
	}

	fmt.Println(messageId)

	return q.redisClient.Get(q.ctx, fmt.Sprintf("%s:active:%s", q.tenantId, messageId)).Result()

	// return q.redisClient.RPopLPush(q.ctx, q.buildListKey(topic, TODO), q.buildListKey(topic, IN_PROGRESS)).Result()
}

// func (q *Queue) Acknowledge(tenant, topic string, messageId string) error {
// 	return q.redisClient.LMove(q.ctx, q.buildListKey(topic, IN_PROGRESS), q.buildListKey(topic, DONE), "right", "right").Err()
// }
