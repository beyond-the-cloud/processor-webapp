package prom

import "github.com/prometheus/client_golang/prometheus"

var (
	HackerCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:      "hackernews_count",
			Help:      "visit HackerNews counter",
		})

	HelloCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:      "hello_counter",
			Help:      "visit /hello endpoint counter",
		})

	StoryCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:      "story_counter",
			Help:      "visit /story endpoint counter",
		})

	GetStoryDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:      "get_story_duration_histogram",
			Help:      "get story duration histogram",
		})

	QueryStoryDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:      "query_story_duration_histogram",
			Help:      "query story duration histogram",
		})

	CreateStoryDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:      "create_story_duration_histogram",
			Help:      "create story duration histogram",
		})

	ConsumeKafkaDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:      "consume_kafka_duration_histogram",
			Help:      "consume kafka duration histogram",
		})

	ElasticSearchDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:      "elastic_search_duration_histogram",
			Help:      "elastic search add story duration histogram",
		})
)
