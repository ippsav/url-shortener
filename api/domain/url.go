package domain

import (
	"net/url"
	"time"
)

type Url struct {
	ID         int64
	Name       string
	RedirectTo string
	OwnerID    string
	CreatedAt  time.Time
}

func (u *Url) Validate() bool {
	_, err := url.ParseRequestURI(u.RedirectTo)
	if err != nil {
		return false
	}
	return true
}
