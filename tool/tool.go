package tool

import (
	"encoding/json"
	"fmt"
	elasticsearch "github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"net/http"
	"processor-webapp/entity"
)

// GetStoryFromHackerNews gets story information from HackerNews
func GetStoryFromHackerNews(id int) (entity.Story, error) {
	url := fmt.Sprintf("%s/%v.%s", "https://hacker-news.firebaseio.com/v0/item", id, "json?print=pretty")
	resp, err := http.Get(url)
	if err != nil {
		return entity.Story{}, err
	}
	defer resp.Body.Close()

	var story entity.Story
	if err := json.NewDecoder(resp.Body).Decode(&story); err != nil {
		return entity.Story{}, err
	}
	return story, nil
}

func GetMaxId() (int, error) {
	url := "https://hacker-news.firebaseio.com/v0/maxitem.json?print=pretty"
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var id int
	if err := json.NewDecoder(resp.Body).Decode(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// InitKafkaConsumer connects to Kafka cluster and subscribe to topic
func InitKafkaConsumer(topic string, server string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"group.id":          "cg",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	consumer.SubscribeTopics([]string{topic}, nil)

	return consumer, nil
}

// ConsumeMsg gets messages from kafka
func ConsumeMsg(c *kafka.Consumer) string {
	msg, err := c.ReadMessage(-1)
	if err == nil {
		return string(msg.Value)
	} else {
		// The client will automatically try to recover from all errors.
		log.Errorf("Consumer error: %v (%v)\n", err, msg)
		return ""
	}
}

// InitElasticSearch initializes ElasticSearch Client
func InitElasticSearch(ESAddr string) (*elasticsearch.Client, error) {
	esClient, err := elasticsearch.NewClient(elasticsearch.SetURL(ESAddr),
		elasticsearch.SetSniff(false),
		elasticsearch.SetHealthcheck(false))
	if err != nil {
		return nil, err
	}
	log.Infof("initialized elasticsearch %v", elasticsearch.Version)
	return esClient, nil
}