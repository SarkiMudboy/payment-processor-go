package payments

import (
	"errors"
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

type card struct {
	Id      string    `json:"id"`
	User    user      `json:"user"`
	Issuer  string    `json:"issuer"`
	Number  string    `json:"number"`
	Expiry  time.Time `json:"expiry"`
	CVV     string    `json:"cvv"`
	Balance float64   `json:"balance,omitempty"`
	Limit   float64   `json:"-"`
}

func (c *card) Charge(amount float64) error {

	if !c.Expired() {
		if c.Balance > amount && amount < c.Limit {
			c.Balance -= amount
			return nil
		}

		return insufficientError
	}

	return errors.New("Card expired!")
}

func (c *card) Credit(amount float64) error {
	c.Balance += amount
	return nil
}

func (c *card) Expired() bool {
	if time.Now().UTC().After(c.Expiry) {
		return true
	}
	return false
}

func issueInvoice(t *transaction, p string) {

	date := time.Now().Format(time.RFC3339)

	i := `
	-----------(%s)---------------
	Name: %s
	Transaction: %s
	Invoice Number: %s
	------------------------------
	Amount: %f
	Tax: %0.00
	Total: %f
	Date: %s
	------------------------------
	Transaction status: %s
	Confirmation code: %s
	`
	i = fmt.Sprintf(i, p, t.User.FullName, t.Id, NewUUID(), t.Amount, t.Amount,
		date, t.Status, t.ConfirmationCode)

	fmt.Println(i)

	t.Invoice = i

}
