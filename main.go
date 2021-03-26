package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"processor-webapp/config"
	"processor-webapp/controller"
	"processor-webapp/entity"
	"processor-webapp/prom"
	"processor-webapp/tool"
	"regexp"
	"strconv"
)

func main() {
	// set log output to stdout
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// connect to database
	dsn := config.DbURL(config.BuildDBConfig())
	var err error
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(err)
	}
	config.DB.AutoMigrate(&entity.Story{})

	// init prom client
	go tool.StartPromClient()

	// register metrics
	prometheus.MustRegister(prom.HelloCounter)
	prometheus.MustRegister(prom.StoryCounter)
	prometheus.MustRegister(prom.GetStoryDuration)
	prometheus.MustRegister(prom.HackerCounter)

	// run router to provide liveness and readiness test
	go tool.InitRouter().Run()

	// initialize kafka consumer
	topic := os.Getenv("DBSchema")
	log.Infof("topic: %v", topic)
	server := os.Getenv("KafkaServer")
	log.Infof("Kafka bootstrap server: %v", server)
	consumer, err := tool.InitKafkaConsumer(topic, server)
	if err != nil {
		log.Error(err)
	} else {
		log.Info("initializing kafka consumer succeed")
	}
	defer consumer.Close()

	// get the max id
	maxId, err := tool.GetMaxId()
	if err != nil {
		log.Error(err)
	}
	log.Infof("got max id: %v from HackerNews", maxId)

	// initialize elasticsearch client
	ESAddr := os.Getenv("ESAddr")
	es, err := tool.InitElasticSearch(ESAddr)
	if err != nil {
		log.Error("error initializing elastic search client: %s", err)
	}
	ctx := context.Background()

	// keep receiving messages from kafka
	for {
		// get messages from kafka consumer
		msg := tool.ConsumeMsg(consumer)
		if msg != "" {
			log.Infof("kafka consumer got message: %v from topic: %v", msg, topic)
			// check if id is valid
			match, err := regexp.MatchString(`^[0-9]*$`, msg)
			if err != nil {
				log.Error(err)
				continue
			}
			if !match {
				log.Errorf("%v is invalid id", msg)
				continue
			}

			// convert string into int
			id, err := strconv.Atoi(msg)
			if err != nil {
				log.Error(err)
			}

			// check if id is smaller than or equal to max id
			if maxId != 0 && id > maxId {
				log.Debugf("id %v is larger than maxID %v", id, maxId)
				continue
			}
			log.Infof("%v is valid id", id)

			// if the story doesn't exist in db
			if !controller.QueryStoryByID(id) {
				// get data from hankernews
				var story entity.Story
				story, err = tool.GetStoryFromHackerNews(id)
				if err != nil {
					log.Error(err)
				} else {
					log.Infof("got the story %d information from HackerNews", id)
				}
				if story.Title == "" {
					log.Infof("got %v %v", story.Type, story)
					continue
				}

				// add story into db
				if err = controller.CreateStory(story); err != nil {
					log.Error(err)
				} else {
					log.Infof("added story %d in db", id)
				}

				// store story in elastic search
				index, err := tool.AddStoryInES(story, es, topic, ctx)
				if err != nil {
					log.Error(err)
				} else {
					log.Infof("index: %v", index)
					log.Infof("added story %d in elastic search", id)
				}
			} else {
				log.Infof("story %d already exists", id)
			}
		} else {
			log.Debug("got blank string")
		}
	}
}
