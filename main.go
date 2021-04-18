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
	"strings"
	"time"
)

func main() {
	// set log output to stdout
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// connect to database
	dsn := config.DbURL(config.BuildDBConfig())
	var err error
	config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// if meet errors connecting to db, retry 3 times
	count := 20
	for err != nil && count>0 {
		log.Debugf("got %v retrying connecting %v times to db", err, count)
		time.Sleep(3000 * time.Millisecond)
		config.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		count--
	}
	if err != nil {
		log.Error(err)
	}

	if err = config.DB.AutoMigrate(&entity.Story{}); err != nil {
		log.Error(err)
	}

	// get all tables
	tables := []string{}
	config.DB.Select(&tables, "SHOW TABLES")
	log.Infof("got %v tables: %v", len(tables), tables)

	// check if table stories exists
	isExist := false
	if len(tables) != 0 {
		for i, table := range tables {
			log.Infof("table %v: %v", i, table)
			if strings.Compare(table, "stories") == 0 {
				isExist = true
				log.Infof("table stories exists")
				break
			}
		}
	}
	log.Infof("table stories exists: %v", isExist)

	// create table stories if not exist
	if len(tables) == 0 || !isExist {
		config.DB.Exec("CREATE TABLE `stories` (\n  `id` bigint(20) NOT NULL AUTO_INCREMENT,\n  `author` longtext,\n  `descendants` bigint(20) DEFAULT NULL,\n  `score` bigint(20) DEFAULT NULL,\n  `create_date` bigint(20) DEFAULT NULL,\n  `title` longtext,\n  `type` longtext,\n  `url` longtext,\n  PRIMARY KEY (`id`)\n) ENGINE=InnoDB AUTO_INCREMENT=26469754 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;")
		log.Info("creating stories table failed, created stories table manually")
		// check tables again
		tables = []string{}
		config.DB.Select(&tables, "SHOW TABLES")
		log.Infof("got %v tables: %v", len(tables), tables)
	}

	// register metrics
	prometheus.MustRegister(prom.HelloCounter)
	prometheus.MustRegister(prom.StoryCounter)
	prometheus.MustRegister(prom.GetStoryDuration)
	prometheus.MustRegister(prom.HackerCounter)
	prometheus.MustRegister(prom.QueryStoryDuration)
	prometheus.MustRegister(prom.CreateStoryDuration)
	prometheus.MustRegister(prom.ConsumeKafkaDuration)
	prometheus.MustRegister(prom.ElasticSearchDuration)

	// run router to provide liveness and readiness test
	go tool.InitRouter()

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
