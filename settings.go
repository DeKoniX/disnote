package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Setting struct {
	Discord struct {
		Token     string
		ChannelID string
	}
}

func Settings() (setting Setting) {
	setting = Setting{}

	dat, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Println("Нет файла настроек: ", err)
	}

	err = yaml.Unmarshal(dat, &setting)
	if err != nil {
		log.Println("Не могу прочитать файл настроек: ", err)
	}

	return setting
}
