package redis

import (
	"fmt"

	redis "gopkg.in/redis.v4"
)

type Message struct {
	noteId    int
	messageId string
}

type Redis struct {
	Options struct {
		Address: string,
		Password: string,
	}
	Client *redis.Client
}

var Client Redis

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
}

// func RedisClient(addr, password string) *redis.Client {
// 	client := redis.NewClient(&redis.Options{
// 		Addr:     addr,
// 		Password: password,
// 		DB:       0,
// 	})

// 	return client
// }

func (client *redis.Client) SetMessage(message Message) (err error) {
	err = client.HSet("DisnoteMessages", fmt.Sprintf("%b", message.noteId), message.messageId).Err()

	return err
}

func (client *redis.Client) GetMessage(noteId int) (messageId []interface{}, err error) {
	strNoteId := fmt.Sprintf("%b", noteId)
	messageId, err = client.HMGet("DisnoteMessages", strNoteId).Result()

	return messageId, err
}

func (client *redis.Client) DelMessage(noteId int) (err error) {
	strNoteId := fmt.Sprintf("%b", noteId)
	err = client.HDel("DisnoteMessages", strNoteId).Err()

	return err
}
