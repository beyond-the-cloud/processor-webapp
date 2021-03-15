package tool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"processor-webapp/entity"
)

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
