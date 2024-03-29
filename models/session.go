package models

import "time"

type Session struct {
	Username      string
	Session_token string
	Expiry        time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
