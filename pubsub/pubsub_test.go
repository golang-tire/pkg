package pubsub

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/golang-tire/pkg/pubsub/test"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"

	"github.com/alicebob/miniredis/v2"
)

var (
	redisClient *redis.Client
	ctx         context.Context
)

func TestMain(m *testing.M) {
	var err error
	redisServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	ctx = context.Background()

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisServer.Host(), redisServer.Port()),
		Password: "",
	})

	code := m.Run()
	redisServer.Close()
	os.Exit(code)
}

func TestNew(t *testing.T) {
	service := New(redisClient)
	assert.NotNil(t, service)
}

func Test_getOrCreateTopic(t *testing.T) {

	service := &service{
		redisClient: redisClient,
		topics:      make(map[string]*Topic),
	}

	// when we dont have the topic already
	topic := service.getOrCreateTopic("test-topic")
	assert.NotNil(t, topic)
	assert.Nil(t, topic.handlers)

	tp, ok := service.topics["test-topic"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, tp)

	//
	tp1, ok := service.topics["new-test-topic"]
	assert.Equal(t, false, ok)
	assert.Nil(t, tp1)

	tp2 := service.getOrCreateTopic("new-test-topic")
	_, ok = service.topics["new-test-topic"]
	assert.Equal(t, true, ok)
	assert.NotNil(t, tp2)
}

func TestPubSub(t *testing.T) {
	service := &service{
		redisClient: redisClient,
		topics:      make(map[string]*Topic),
	}

	// should be zero when start
	num, err := redisClient.PubSubNumPat(ctx).Result()
	assert.Nil(t, err)
	assert.Equal(t, int64(0), num)

	var wait = make(chan bool, 1)
	var receivedData test.HelloWorld

	service.Subscribe(ctx, "test-channel", func(ctx context.Context, msg *test.HelloWorld) {
		receivedData.Name = msg.Name
		wait <- true
	})

	err = service.Publish(ctx, "test-channel", &test.HelloWorld{
		Name: "test-one",
	})
	assert.Nil(t, err)

	<-wait
	assert.Equal(t, "test-one", receivedData.Name)
}
