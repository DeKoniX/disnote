package main

import (
	"fmt"

	redis "gopkg.in/redis.v4"
)

type Message struct {
	noteId    int
	messageId string
}

func RedisClient(addr, password string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return client
}

func RedisSetMessage(client *redis.Client, message Message) (err error) {
	err = client.HSet("DisnoteMessages", fmt.Sprintf("%b", message.noteId), message.messageId).Err()

	return err
}

func RedisGetMessage(client *redis.Client, noteId int) (messageId []interface{}, err error) {
	strNoteId := fmt.Sprintf("%b", noteId)
	messageId, err = client.HMGet("DisnoteMessages", strNoteId).Result()

	return messageId, err
}

func RedisDelMessage(client *redis.Client, noteId int) (err error) {
	strNoteId := fmt.Sprintf("%b", noteId)
	err = client.HDel("DisnoteMessages", strNoteId).Err()

	return err
}
