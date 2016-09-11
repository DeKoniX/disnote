package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	redis "gopkg.in/redis.v4"

	"github.com/bwmarrin/discordgo"
)

var BotID string
var ChannelID string
var DB *sql.DB
var RClient *redis.Client

func main() {
	setting := Settings()
	ChannelID = setting.Discord.ChannelID

	RClient = RedisClient(setting.Redis.Address, setting.Redis.Password)

	DB = db_init()
	go run_bot(setting.Discord.Token)

	log.Println("Бот запущен")

	<-make(chan struct{})
	return
}

func run_bot(token string) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Не могу провести авторизацию: ", err)
		return
	}

	u, err := dg.User("@me")
	if err != nil {
		log.Println("Ошибка получение аккаунта: ", err)
	}

	BotID = u.ID

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Println("Ошибка присоединения: ", err)
		return
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	_ = s.UpdateStatus(0, "DisNote - ThisIsNotes!")

	clearChannel(s)
	postAllMess(s)
}

func clearChannel(s *discordgo.Session) {
	var mass_message_id []string

	mass_message, err := s.ChannelMessages(ChannelID, 100, "", "")

	if err != nil {
		log.Println(err)
	}

	for len(mass_message) != 0 {
		for _, message := range mass_message {
			mass_message_id = append(mass_message_id, message.ID)
		}
		err := s.ChannelMessagesBulkDelete(ChannelID, mass_message_id)
		if err != nil {
			log.Println(err)
		}

		mass_message, err = s.ChannelMessages(ChannelID, 100, "", "")
		if err != nil {
			log.Println(err)
		}
	}
}

func postAllMess(s *discordgo.Session) {
	rows := db_select(DB)
	for _, row := range rows {
		str := fmt.Sprintf("%d -> %s (<@%s>)", row.id, row.text, row.user_id)
		mess, _ := s.ChannelMessageSend(ChannelID, str)
		err := RedisSetMessage(RClient, Message{
			noteId:    int(row.id),
			messageId: mess.ID,
		})
		if err != nil {
			log.Panic(err)
		}
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if ChannelID != m.ChannelID {
		return
	}

	if strings.HasPrefix(m.Content, "-add") {
		id := db_insert(DB, strings.TrimPrefix(m.Content, "-add"), m.Author.ID)

		str := fmt.Sprintf("%d -> %s (<@%s>)", id, strings.TrimPrefix(m.Content, "-add"), m.Author.ID)
		mess, _ := s.ChannelMessageSend(ChannelID, str)

		err := RedisSetMessage(RClient, Message{
			noteId:    int(id),
			messageId: mess.ID,
		})
		if err != nil {
			log.Panic(err)
		}
	}

	if strings.HasPrefix(m.Content, "-del ") {
		id, err := strconv.Atoi(strings.TrimPrefix(m.Content, "-del "))
		if err != nil {
			send_sleep_and_del("Введите пожалуйста число -del *num*", s)
		} else {
			if db_delete(DB, id) == true {
				mess, err := RedisGetMessage(RClient, id)
				if err != nil {
					log.Panic(err)
				}
				IDMess := fmt.Sprintf("%v", mess[0])
				_ = s.ChannelMessageDelete(ChannelID, IDMess)
				err = RedisDelMessage(RClient, id)
				if err != nil {
					log.Panic(err)
				}
				str := fmt.Sprintf("Заметка %v удалена", id)
				send_sleep_and_del(str, s)
			} else {
				str := fmt.Sprintf("Заметка %v не существует", id)
				send_sleep_and_del(str, s)
			}
		}
	}

	if m.Content == "-clear" {
		clearChannel(s)
		postAllMess(s)
	}

	if m.Content == "-help" {
		send_sleep_and_del("\n-add - добавить заметку\n-del <num> - удалить заметку\n-clear - очистить канал и заного написать заметки", s)
	}

	_ = s.ChannelMessageDelete(ChannelID, m.ID)
}

func send_sleep_and_del(message string, s *discordgo.Session) bool {
	mess, _ := s.ChannelMessageSend(ChannelID, message)
	time.Sleep(time.Second * 10)
	_ = s.ChannelMessageDelete(ChannelID, mess.ID)
	return true
}
