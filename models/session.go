package models

import "time"

type Session struct {
	Username      string
	Session_token string
	Expiry        time.Time
}

func (s *Session) Is_expired() bool {
	return false
	//TODO FIX
	//return s.Expiry.Before(time.Now())
}
