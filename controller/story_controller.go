package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math/rand"
	"net/http"
	"processor-webapp/entity"
	"processor-webapp/model"
)

// GetAllStories ... Get all stories
func GetAllStories(c *gin.Context) {
	var stories []entity.Story
	err := model.GetAllStories(stories)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, story := range stories {
		c.JSON(http.StatusOK, gin.H{
			"id":          story.ID,
			"by":          story.Author,
			"descendants": story.Descendants,
			"score":       story.Score,
			"time":        story.CreateDate,
			"title":       story.Title,
			"type":        story.Type,
			"url":         story.URL,
		})
	}
}

// GetStories ... Get all stories
func GetStories() ([]entity.Story, error) {
	var stories []entity.Story
	err := model.GetAllStories(stories)
	if err != nil {
		return nil, err
	}
	return stories, nil
}

// CreateStory ... Create Story
func CreateStory(story entity.Story) error {
	err := model.CreateStory(story)
	if err != nil {
		return err
	}
	return nil
}

// QueryStoryByIDGin ... Checks connection to database
func QueryStoryByIDRouter(c *gin.Context) {
	// randomly get an id in the range of (10000, 26497625]
	min := 100000
	max := 26497625
	id := rand.Intn(max-min)+min
	var story entity.Story
	err := model.GetStoryByID(story, id)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		log.Errorf("visiting /story got error: %v", err)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg": "ready",
		})
		log.Info("visiting /story succeed")
	}
}

// QueryStoryByID ... Query the story by id
func QueryStoryByID(id int) bool {
	return model.QueryStoryByID(id)
}

// GetStoryByID ... Get the story by id
func GetStoryByID(id int) (entity.Story, error) {
	var story entity.Story
	err := model.GetStoryByID(story, id)
	if err != nil {
		return story, err
	}
	return story, nil
}

// GetStoryByTitle ... Get the story by title
func GetStoryByTitle(title string) (entity.Story, error) {
	var story entity.Story
	err := model.GetStoryByTitle(story, title)
	if err != nil {
		return story, err
	}
	return story, nil
}

// UpdateStory ... Update the story information
func UpdateStory(story entity.Story) error {
	if err := model.UpdateStory(story, story.ID); err != nil {
		return err
	}
	return nil
}

// DeleteStory ... Delete the story
func DeleteStory(id int) error {
	var story entity.Story
	err := model.DeleteStory(story, id)

	if err != nil {
		return err
	}
	return nil
}
