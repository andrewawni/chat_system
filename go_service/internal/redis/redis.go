package redis

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

// Client - redis client wrapper
type Client struct {
	rdb *redis.Client
	ctx context.Context
}

type transactionFunction func(tx *redis.Tx) error

// CreateClient - create redis client
func CreateClient(connectionURL string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     connectionURL,
		Password: "",
		DB:       0,
	})

	if rdb == nil {
		return &Client{}, errors.New("Invalid URL")
	}

	client := Client{
		rdb: rdb,
		ctx: context.Background(),
	}
	log.Printf("[redis] client connected successfully")
	return &client, nil
}

// SetInt - Set key to val
func (client *Client) SetInt(key string, val int) error {
	err := client.rdb.Set(client.ctx, key, val, 0).Err()
	return err
}

// GetInt - Get value of key
func (client *Client) GetInt(key string) (int, error) {
	val, err := client.rdb.Get(client.ctx, key).Result()
	if err != nil {
		return -1, err
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return intVal, err
}

// IncrementIntAndSetNewKey - increments the value of key from n to n + 1, and sets a new redis key = `key:n+1` and value = val
func (client *Client) IncrementIntAndSetNewKey(key string, val int) (int, error) {
	num := -1
	txf := func(tx *redis.Tx) error {
		// Get current value, fail if key doesn't exist (returns redis.Nil)
		n, err := tx.Get(client.ctx, key).Int()
		if err != nil {
			return err
		}
		// Actual operation (local in optimistic lock).
		n++
		// Operation is committed only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(client.ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(client.ctx, key, n, 0)
			// set new counter to 0
			newKey := fmt.Sprint(key, ":", n)
			pipe.Set(client.ctx, newKey, val, 0)
			num = n
			return nil
		})
		return err
	}
	err := client.commitTransaction(key, txf)
	return num, err
}

func (client *Client) commitTransaction(key string, txf transactionFunction) error {
	const maxRetries int = 1000
	for i := 0; i < maxRetries; i++ {
		err := client.rdb.Watch(client.ctx, txf, key)
		if err == nil {
			// Success.
			log.Printf("[redis] transaction committed successfully after %d retries", i)
			return nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return err
	}
	return errors.New("[redis] reached maximum number of retries")
}
