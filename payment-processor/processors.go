package main

import (
	"fmt"
	"time"
)

type Processor interface {
	oneTimePayment()
	subscriptionPayment()
	refund()
	invoice()
}

type Subscription struct {
	user   *User
	name   string
	period string
	due    time.Time
	active bool
	plan   float64
	billed float64
}

// func (s *Subscription) VerifyBilling() bool {

// }

// func (s Subscription) Biil() bool {

// }

func processor() {
	fmt.Println("payment")
}
