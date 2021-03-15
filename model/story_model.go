package model

import (
	"errors"
	"processor-webapp/config"
	"processor-webapp/entity"
)

// GetAllStories Fetch all story data
func GetAllStories(stories []entity.Story) (err error) {
	if err = config.DB.Find(&stories).Error; err != nil {
		return err
	}
	return nil
}

// CreateStory ... Insert New data
func CreateStory(story entity.Story) (err error) {
	if err = config.DB.Create(&story).Error; err != nil {
		return err
	}
	return nil
}

// QueryStoryByID ... Query story by Id
func QueryStoryByID(id int) bool {
	var story entity.Story
	result := config.DB.Where("id = ?", id).First(&story)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

// GetStoryByID ... Fetch story by Id
func GetStoryByID(story entity.Story, id int) (err error) {
	if err = config.DB.Where("id = ?", id).First(&story).Error; err != nil {
		return err
	}
	return nil
}

// GetStoryByTitle ... Fetch story by Title
func GetStoryByTitle(story entity.Story, title string) (err error) {
	if err = config.DB.Where("title = ?", title).First(&story).Error; err != nil {
		return err
	}
	return nil
}

// UpdateStory ... Update story
// Not useful in this case, but still keep it
func UpdateStory(story entity.Story, id int) (err error) {
	var originStory entity.Story
	if err = config.DB.Where("id = ?", id).First(&originStory).Error; err != nil {
		return err
	}

	config.DB.Save(story)
	return nil
}

// DeleteStory ... Delete story
func DeleteStory(story entity.Story, id int) (err error) {
	var originStory *entity.Story
	if err = config.DB.Where("id = ?", id).First(&originStory).Error; err != nil {
		return err
	}

	if originStory == nil {
		return errors.New("story with this id doesn't exist")
	}

	config.DB.Where("id = ?", id).Delete(&story)
	return nil
}
