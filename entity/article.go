package entity

import (
	"time"
)

type Article struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Article) FromMap(m map[string]any) *Article {
	a.ID = m["id"].(string)
	a.URL = m["url"].(string)
	a.Title = m["title"].(string)
	a.CreatedAt = m["created_at"].(time.Time)
	return a
}
