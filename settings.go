package main

import (
	"io/ioutil"
	"log"
	"os/user"

	"gopkg.in/yaml.v2"
)

type Setting struct {
	Discord struct {
		Token     string
		ChannelID string
	}
}

func Settings() (setting Setting) {
	setting = Setting{}
	usr, _ := user.Current()

	dat, err := ioutil.ReadFile(usr.HomeDir + "/.config/disnote.yml")
	if err != nil {
		log.Panicln("Нет файла настроек: ", err)
	}

	err = yaml.Unmarshal(dat, &setting)
	if err != nil {
		log.Println("Не могу прочитать файл настроек: ", err)
	}

	return setting
}
