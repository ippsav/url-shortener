package domain

import (
	"net/url"
	"strings"
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
	chkUrl, err := url.ParseRequestURI(u.RedirectTo)
	if err != nil || strings.Contains(chkUrl.Host, ".") {
		return false
	}
	return true
}
