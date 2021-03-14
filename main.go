package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"processor-webapp/config"
	"processor-webapp/controller"
	"processor-webapp/tool"
)

func main() {
	dsn := config.DbURL(config.BuildDBConfig())
	var err error
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
	}

	Topic := os.Getenv("DBSchema")
	fmt.Println(Topic)

	// set up kafka consumer

	// get ids
	id := 26422799

	// if the story already exists in db
	story, err := controller.GetStoryByID(id)
	if err != nil {
		log.Error(err)
	}
	if story != nil {
		log.Info("The story with id: %v already exists in the database", id)
	}

	// if the story doesn't exist in db, get data from hankernews
	story, err = tool.GetStoryFromHackerNews(id)
	if err != nil {
		log.Error(err)
	}
	// add story into db
	if err = controller.CreateStory(story); err != nil {
		log.Error(err)
	}

	// store story in elastic search
}
