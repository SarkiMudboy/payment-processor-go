package main

import "time"

type Subscription struct {
	user   *User
	name   string
	period string
	due    time.Time // in UTC
	active bool
	plan   float64
	billed float64
}

func (s *Subscription) VerifyBilling() bool {

	now := time.Now().UTC()

	if (now.Equal(s.due) || now.After(s.due)) && s.active {
		return true
	}

	return false
}

func (s *Subscription) Biil() bool {

	now := time.Now().UTC()
	s.billed += s.plan
	s.due = now.Add(Period[s.period])

	return true

}
