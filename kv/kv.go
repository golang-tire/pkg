package kv

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-tire/pkg/log"
)

var (
	client *Client
)

// Config kv redis config
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Client kv client service
type Client struct {
	ctx    context.Context
	client *redis.Client
}

// With returns a Client with context
func (c *Client) With(ctx context.Context) *redis.Client {
	c.ctx = ctx
	return c.client.WithContext(ctx)
}

// Get will return redis connection
func Get() *Client {
	return client
}

// Conn will return redis connection
func (c *Client) Conn() *redis.Conn {
	return c.client.Conn(c.ctx)
}

// Set will set a key ( string ) to a value ( interface )
func (c *Client) Set(key string, val interface{}, alive time.Duration) error {
	err := c.client.Set(c.ctx, key, val, alive).Err()
	if err != nil {
		log.Error("kv: write error", log.Err(err))
	}
	return err
}

// GetString will read and return a key as string
func (c *Client) GetString(key string) (string, error) {
	val, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, err
}

// Delete will remove a key
func (c *Client) Delete(key string) error {
	err := c.client.Del(c.ctx, key).Err()
	if err != nil {
		log.Error("kv: delete key error", log.Err(err))
		return err
	}
	return nil
}

// Init will initialize the key value store
func Init(ctx context.Context, config *Config) (*Client, error) {

	// use default config
	if config == nil {
		config = &Config{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Error("kv: connect error", log.Err(err))
		return nil, err
	}

	go func() {
		<-ctx.Done()
		if err := rdb.Close(); err != nil {
			log.Error("error in close redis connection", log.Err(err))
		}
	}()

	client = &Client{ctx: ctx, client: rdb}
	return client, nil
}

// InitWithConn initialize the key value store using a given redis client
func InitWithConn(ctx context.Context, redisClient *redis.Client) (*Client, error) {

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Error("kv: connect error", log.Err(err))
		return nil, err
	}

	go func() {
		<-ctx.Done()
		if err := redisClient.Close(); err != nil {
			log.Error("error in close redis connection", log.Err(err))
		}
	}()

	client = &Client{ctx: ctx, client: redisClient}
	return client, nil
}
