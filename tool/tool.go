package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	elasticsearch "github.com/olivere/elastic/v7"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"net/http"
	"processor-webapp/controller"
	"processor-webapp/entity"
	"processor-webapp/prom"
	"time"
)

// GetStoryFromHackerNews gets story information from HackerNews
func GetStoryFromHackerNews(id int) (entity.Story, error) {
	url := fmt.Sprintf("%s/%v.%s", "https://hacker-news.firebaseio.com/v0/item", id, "json?print=pretty")
	resp, err := http.Get(url)
	if err != nil {
		return entity.Story{}, err
	}
	defer resp.Body.Close()

	prom.HackerCounter.Inc()

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

	prom.HackerCounter.Inc()

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
	defer func(begun time.Time) {
		prom.ConsumeKafkaDuration.Observe(time.Since(begun).Seconds())
	}(time.Now())
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

// InitRouter initializes router and provides Liveness and Readiness Test
func InitRouter() *gin.Engine {
	r := gin.Default()

	// tests if the processor-webapp runs normally
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
		log.Info("visiting /hello succeed")
		prom.HelloCounter.Inc()
	})

	// tests if the processor-webapp connects to database normally
	r.GET("/story", controller.QueryStoryByIDRouter)

	return r
}

// StartPromClient initializes Prometheus Client
func StartPromClient() {
	log.Info("initializing prometheus client")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

// AddStoryInES adds story in ElasticSearch
func AddStoryInES(story entity.Story, es *elasticsearch.Client, topic string, ctx context.Context) (*elasticsearch.IndexResponse, error) {
	esStory := entity.EsStory{
		ID:    story.ID,
		Title: story.Title,
	}
	storyJSON, err := json.Marshal(esStory)

	defer func(begun time.Time) {
		prom.ElasticSearchDuration.Observe(time.Since(begun).Seconds())
	}(time.Now())

	index, err := es.Index().
		Index(topic).
		BodyJson(string(storyJSON)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	return index, nil
}