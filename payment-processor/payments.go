package main

import (
	"fmt"
	"time"
)

type subscription struct {
	id     string
	user   user
	name   string
	period string
	due    time.Time // in UTC
	active bool
	plan   float64
	billed float64
}

func (s *subscription) VerifyBilling() bool {

	now := time.Now().UTC()

	if (now.Equal(s.due) || now.After(s.due)) && s.active {
		return true
	}

	return false
}

func (s *subscription) Biil() bool {

	s.billed += s.plan
	err := s.SetBilling()

	if err != nil {
		return false
	}

	return true

}

func (s *subscription) SetBilling() error {
	// sets the next billing date for the subscription

	now := time.Now().UTC()

	switch s.period {
	case "week":
		s.due = now.Add(week)
	case "month":
		s.due = now.Add(month)
	case "year":
		s.due = now.Add(year)
	default:
		return fmt.Errorf("%s is not a valid billing period", s.period)
	}

	return nil
}

type paypalSubscription struct {
	subscription
}
