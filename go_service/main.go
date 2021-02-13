package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: "",
	DB:       0,
})

type applicationType struct {
	Status           string `json:"status"`
	ApplicationToken string `json:"application_token"`
	ApplicationName  string `json:"application_name"`
}

type chatType struct {
	Status           string `json:"status"`
	ApplicationToken string `json:"application_token"`
	ChatNumber       int    `json:"chat_number"`
	ChatName         string `json:"chat_name"`
}

type messageType struct {
	Status           string `json:"status"`
	ApplicationToken string `json:"application_token"`
	ChatNumber       int    `json:"chat_number"`
	MessageNumber    int    `json:"message_number"`
	MessageContent   string `json:"message_content"`
}

type transactionFunction func(tx *redis.Tx) error

func setKey(key string, val int) error {
	err := rdb.Set(ctx, key, val, 0).Err()
	return err
}

func commitTransaction(key string, txf transactionFunction) error {
	const maxRetries int = 100
	for i := 0; i < maxRetries; i++ {
		err := rdb.Watch(ctx, txf, key)
		if err == nil {
			// Success.
			return nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return err
	}
	return errors.New("increment reached maximum number of retries")
}

func incrementKeyTransaction(key string) (int, error) {
	// Transactional function.
	const maxRetries int = 100
	var value int = -1
	txf := func(tx *redis.Tx) error {
		// Get current value, fail if key doesn't exist (returns redis.Nil)
		n, err := tx.Get(ctx, key).Int()
		if err != nil {
			return err
		}
		// Actual operation (local in optimistic lock).
		n++
		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n, 0)
			value = n
			return nil
		})
		return err
	}

	for i := 0; i < maxRetries; i++ {
		err := rdb.Watch(ctx, txf, key)
		if err == nil {
			// Success.
			return value, nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return value, err
	}

	return value, errors.New("increment reached maximum number of retries")
}

func createApplication(w http.ResponseWriter, r *http.Request) {

	token := uuid.New().String()

	err := setKey(token, 0)

	if err != nil {
		panic(err)
	}

	payload := applicationType{
		Status:           "success",
		ApplicationToken: token,
	}

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(payload)
	w.Write(response)
}

func createChat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["application_token"]
	var num int = -1
	txf := func(tx *redis.Tx) error {
		// Get current value, fail if key doesn't exist (returns redis.Nil)
		n, err := tx.Get(ctx, token).Int()
		if err != nil {
			return err
		}
		// Actual operation (local in optimistic lock).
		n++
		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, token, n, 0)
			num = n
			// set chat counter to 0
			chatKey := fmt.Sprint(token, ":", num)
			pipe.Set(ctx, chatKey, 0, 0)
			return nil
		})
		return err
	}
	err := commitTransaction(token, txf)
	if err != nil {
		panic(err)
	}
	payload := chatType{
		Status:           "success",
		ApplicationToken: token,
		ChatNumber:       num,
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(payload)
	w.Write(response)
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["application_token"]
	chatNumber, err := strconv.Atoi(vars["chat_number"])
	if err != nil {
		panic(err)
	}

	key := fmt.Sprint(token, ":", chatNumber)
	var num int = -1
	txf := func(tx *redis.Tx) error {
		// Get current value, fail if key doesn't exist (returns redis.Nil)
		n, err := tx.Get(ctx, key).Int()
		if err != nil {
			return err
		}
		// Actual operation (local in optimistic lock).
		n++
		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n, 0)
			num = n
			// set chat counter to 0
			messageKey := fmt.Sprint(key, ":", num)
			pipe.Set(ctx, messageKey, 0, 0)
			return nil
		})
		return err
	}
	err = commitTransaction(key, txf)
	if err != nil {
		panic(err)
	}
	payload := messageType{
		Status:           "success",
		ApplicationToken: token,
		ChatNumber:       chatNumber,
		MessageNumber:    num,
	}
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(payload)
	w.Write(response)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/applications", createApplication).Methods("POST")
	r.HandleFunc("/api/applications/{application_token}/chats", createChat).Methods("POST")
	r.HandleFunc("/api/applications/{application_token}/chats/{chat_number}/messages", createMessage).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
}
