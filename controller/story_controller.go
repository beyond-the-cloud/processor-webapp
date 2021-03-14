package controller

import (
	"processor-webapp/entity"
	"processor-webapp/model"
)

//GetStories ... Get all stories
func GetStories() (*[]entity.Story, error) {
	var stories *[]entity.Story
	err := model.GetAllStories(stories)
	if err != nil {
		return nil, err
	}
	return stories, nil
}

//CreateStory ... Create Story
func CreateStory(story *entity.Story) error {
	err := model.CreateStory(story)
	if err != nil {
		return err
	}
	return nil
}

//GetStoryByID ... Get the story by id
func GetStoryByID(id int) (*entity.Story, error) {
	var story *entity.Story
	err := model.GetStoryByID(story, id)
	if err != nil {
		return nil, err
	}
	return story, nil
}

//GetStoryByTitle ... Get the story by title
func GetStoryByTitle(title string) (*entity.Story, error) {
	var story *entity.Story
	err := model.GetStoryByTitle(story, title)
	if err != nil {
		return nil, err
	}
	return story, nil
}

//UpdateStory ... Update the story information
func UpdateStory(story *entity.Story) error {
	if err := model.UpdateStory(story, story.ID); err != nil {
		return err
	}
	return nil
}

//DeleteStory ... Delete the story
func DeleteStory(id int) error {
	var story entity.Story
	err := model.DeleteStory(&story, id)

	if err != nil {
		return err
	}
	return nil
}
