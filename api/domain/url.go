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
	if err != nil {
		return false
	}
	if strings.Contains(chkUrl.Host, ".") && !strings.Contains(chkUrl.Host, "www") {
		return true
	} else if strings.Contains(chkUrl.Host, "www") && strings.Index(chkUrl.Host, ".") != strings.LastIndex(chkUrl.Host, ".") {
		return true
	}
	return false
}
