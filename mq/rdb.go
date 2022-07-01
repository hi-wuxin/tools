package mq

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"log"
	"runtime/debug"
	"sync"
)

type RMQ struct {
	ctx      context.Context
	cmd      *redis.Client
	mutex    *sync.RWMutex
	handlers map[string]Handler
}

func NewRedisMQ(ctx context.Context, client *redis.Client) MQueue {
	return &RMQ{
		ctx:      ctx,
		cmd:      client,
		mutex:    &sync.RWMutex{},
		handlers: make(map[string]Handler, 0),
	}
}

func (R *RMQ) Register(topics ...string) error {

	for _, topic := range topics {
		R.mutex.RLock()
		R.handlers[topic] = nil
		R.mutex.RUnlock()
	}
	return nil
}

// Publish 发布消息
func (R *RMQ) Publish(topic string, msgs ...interface{}) error {
	//必须注册了才能发布
	//if _, ok := R.handlers[topic]; !ok {
	//	return errors.New("topic not register")
	//}
	data := make([]interface{}, len(msgs))
	for _, msg := range msgs {
		body, _ := json.Marshal(Body{Raw: msg})
		s := base64.StdEncoding.EncodeToString(body)
		data = append(data, s)
	}

	if err := R.cmd.RPush(R.ctx, topic, data...).Err(); err != nil {
		return err
	}
	//通知队列去处理
	return R.cmd.Publish(R.ctx, topic, topic).Err()
}

// Subscribe 订阅消息
func (R *RMQ) Subscribe(topic string, handler Handler) error {

	//必须注册了才能发布
	if _, ok := R.handlers[topic]; !ok {
		return errors.New("topic not register")
	}
	if R.handlers[topic] != nil {
		//已被注册过
		return errors.New("topic already register")
	}
	R.mutex.RLock()
	R.handlers[topic] = handler
	R.mutex.RUnlock()
	return nil
}
func (R *RMQ) daemonServer() error {
	for topic, handler := range R.handlers {
		if handler != nil {
			go R.work(topic, handler)
		}
	}
	return nil
}
func (R *RMQ) work(topic string, handler Handler) {

	defer func() {
		if err := recover(); err != nil {
			log.Println("mq work panic:", string(debug.Stack()))
		}
	}()

	sub := R.cmd.PSubscribe(R.ctx, topic)
	sub.Channel()
	for {
		select {
		//终止处理
		case <-R.ctx.Done():
			return
		case <-sub.Channel():
			for R.cmd.LLen(R.ctx, topic).Val() > 0 {
				if msg := R.cmd.LPop(R.ctx, topic).Val(); msg != "" {
					data, err := base64.StdEncoding.DecodeString(msg)
					if err != nil {
						continue
					}
					var body Body
					err = json.Unmarshal(data, &body)
					//TODO 这里可以用协程池来处理
					handler(body.Raw)
				}
			}
		}
	}
}

func (R *RMQ) Length(topic string) (int64, error) {
	return R.cmd.LLen(R.ctx, topic).Result()
}
