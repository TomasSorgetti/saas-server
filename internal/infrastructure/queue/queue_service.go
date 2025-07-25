package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Queue struct {
    client *redis.Client
    name   string
}

func NewQueue(client *redis.Client, name string) *Queue {
    return &Queue{
        client: client,
        name:   name,
    }
}

func (q *Queue) Enqueue(ctx context.Context, job interface{}) error {
    data, err := json.Marshal(job)
    if err != nil {
        return err
    }
    return q.client.LPush(ctx, q.name, data).Err()
}

func (q *Queue) Dequeue(ctx context.Context) ([]byte, error) {
    result, err := q.client.BLPop(ctx, 0*time.Second, q.name).Result()
    if err != nil {
        return nil, err
    }
    if len(result) < 2 {
        return nil, nil
    }
    return []byte(result[1]), nil
}