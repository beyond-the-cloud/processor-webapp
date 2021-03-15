package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"processor-webapp/config"
	"processor-webapp/controller"
	"processor-webapp/entity"
	"processor-webapp/tool"
)

func main() {
	dsn := config.DbURL(config.BuildDBConfig())
	var err error
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
	}
	config.DB.AutoMigrate(&entity.Story{})

	Topic := os.Getenv("DBSchema")
	fmt.Println(Topic)

	// set up kafka consumer

	// get ids
	id := 26460390

	// if the story doesn't exist in db
	if !controller.QueryStoryByID(id) {
		// get data from hankernews
		var story entity.Story
		story, err = tool.GetStoryFromHackerNews(id)
		if err != nil {
			log.Error(err)
		} else {
			log.Infof("Got the story %d information from HackerNews", id)
		}

		// add story into db
		if err = controller.CreateStory(story); err != nil {
			log.Error(err)
		} else {
			log.Infof("Added story %d in db", id)
		}

		// store story in elastic search
		log.Infof("Added story %d in elastic search", id)
	} else {
		log.Infof("Story %d already exists", id)
	}
}