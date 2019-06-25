package model

// TwitterCard response
type TwitterCard struct {
	Card        string `json:"card"`
	Site        string `json:"site"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Retry       bool   `json:"retry"`
}

const summaryCard = "summary"

// ValuePresent is
func (t *TwitterCard) ValuePresent() bool {
	return t.Title != "" || t.Description != "" || t.Image != ""
}
