package pubsub

import (
	"context"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/golang-tire/pkg/log"

	"github.com/go-redis/redis/v8"
)

var pubSubSrv *service

type service struct {
	lock        sync.RWMutex
	redisClient *redis.Client
	topics      map[string]*Topic
}

// Topic topic holder that contain handlers lists
type Topic struct {
	handlers []Handler
}

// Handler is an interface and will contain target function to run in subscriber
type Handler interface{}

// Service pubsub service
type Service interface {
	Publish(ctx context.Context, topic string, msg proto.Message) error
	Subscribe(ctx context.Context, topic string, handler Handler)
}

// New create a new instance of pubsub service
func New(client *redis.Client) Service {
	pubSubSrv = &service{
		redisClient: client,
		topics:      make(map[string]*Topic),
	}
	return pubSubSrv
}

// Get return pubsub service instance
func Get() Service {
	return pubSubSrv
}

func (s *service) getOrCreateTopic(topic string) *Topic {
	_, ok := s.topics[topic]
	if !ok {
		tp := &Topic{
			handlers: nil,
		}
		s.topics[topic] = tp
	}
	return s.topics[topic]
}

// Publish send a proto message to a topic.
func (s *service) Publish(ctx context.Context, topic string, msg proto.Message) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.getOrCreateTopic(topic)

	m, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return s.redisClient.Publish(ctx, topic, m).Err()
}

// Subscribe subscribe on a topic
// handler function should be of format func(ctx context.Context, msg *proto.Message)
func (s *service) Subscribe(ctx context.Context, topic string, handler Handler) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Reflection is slow, but this is done only once on subscriber setup
	handlerFunc := reflect.TypeOf(handler)
	if handlerFunc.Kind() != reflect.Func {
		panic("pubsub: handler needs to be a func")
	}

	if handlerFunc.NumIn() != 2 {
		panic(`pubsub: handler should be of format
		func(ctx context.Context, msg *proto.Message)
		but didn't receive enough args`)
	}

	if handlerFunc.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		panic(`pubsub: handler should be of format
		func(ctx context.Context, msg *proto.Message)
		but first arg was not context.Context`)
	}

	tp := s.getOrCreateTopic(topic)
	tp.handlers = append(tp.handlers, handler)
	ps := s.redisClient.Subscribe(ctx, topic)

	go func() {
		ch := ps.Channel()
		for msg := range ch {
			if msg.Channel == topic {
				obj := reflect.New(handlerFunc.In(1).Elem()).Interface()
				err := proto.Unmarshal([]byte(msg.Payload), obj.(proto.Message))
				if err != nil {
					log.Error("unmarshal data failed", log.String("topic", topic), log.Err(err))
					continue
				}
				for _, h := range tp.handlers {
					fn := reflect.ValueOf(h)
					go func() {
						fn.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(obj)})
					}()
				}
			}
		}
	}()

	go func() {
		<-ctx.Done()
		err := ps.Close()
		log.Error("close subscriber failed", log.String("topic", topic), log.Err(err))
	}()
}
