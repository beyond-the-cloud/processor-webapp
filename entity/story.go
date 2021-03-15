package entity

type Story struct {
	ID          int    `json:"id"` // The item's unique id.
	Author      string `json:"by"` // The username of the item's author.
	Descendants int    `json:"descendants"`
	Score       int    `json:"score"` // The story's score, or the votes for a pollopt.
	CreateDate  int    `json:"time"`  // Creation date of the item, in Unix Time.
	Title       string `json:"title"` // The title of the story, poll or job. HTML.
	Type        string `json:"type"`  // The type of item. One of "job", "story", "comment", "poll", or "pollopt".
	URL         string `json:"url"`   // The URL of the story.
	//Kids        []int  `json:"kids"`
	//Text        string    `json:"text"`  // The comment, story or poll text. HTML.
}