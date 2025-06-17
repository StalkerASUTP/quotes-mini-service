package quote

import "time"

type Quote struct {
	ID        int       `json:"id" example:"1"`
	Author    string    `json:"author" example:"Confucius"`
	Quote     string    `json:"quote" example:"Life is simple, but we insist on making it complicated."`
	CreatedAt time.Time `json:"created_at" example:"2025-05-29T00:00:00Z"`
}

func NewQuote(author, quote string) *Quote {
	return &Quote{
		Author: author,
		Quote:  quote,
	}
}
