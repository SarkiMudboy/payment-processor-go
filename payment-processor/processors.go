package main

import (
	_ "errors"
	"fmt"
	"log"
	"time"
)

type Processor interface {
	oneTimePayment()
	subscriptionPayment()
	refund()
	invoice()
}

func processor() {
	fmt.Println("payment")
}

type CreditCardProcessor struct {
	card  *Card
	Label string
}

func (c *CreditCardProcessor) oneTimePayment(t *Transaction) {

	fmt.Printf("submitting transaction: %s to %s\n", t.Id, c.card.Issuer) // change to logs
	fmt.Println("Verifying transaction...")
	fmt.Println("Transaction approved, processing payment...")

	t.Status = "pending"
	t.ConfirmationCode = NewUUID()

	err := c.card.Charge(t.Amount)

	if err != nil {
		t.Status = "failed"
		log.Fatal("(Error) Could not process payment:", err)
	} else {
		fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
		t.Status = "paid"
	}

	c.Invoice(t)

}

func (c *CreditCardProcessor) Subscription(s *Subscription, t *Transaction) {

	if s.VerifyBilling() {
		fmt.Printf("submitting transaction: %s to %s\n", t.Id, c.card.Issuer) // change to logs
		fmt.Println("Verifying transaction...")
		fmt.Println("Transaction approved, processing payment...")

		t.Status = "pending"
		t.ConfirmationCode = NewUUID()

		err := c.card.Charge(s.plan)

		if err != nil {
			t.Status = "failed"
			log.Fatal("(Error) Could not process payment:", err)
		} else {
			fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
			fmt.Printf("You have renewed your %s subscription plan till %s\n", s.name, s.due)
			t.Status = "paid"
		}
	} else {
		fmt.Println("(Fail) Transaction failed")
	}

	c.Invoice(t)

}

func (c *CreditCardProcessor) Refund(r *Transaction, t *Transaction) {

	entries, err := Load(TransactionFile)

	if err != nil {
		log.Fatal(err)
	}

	r, err = r.Load(entries)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("submitting refund request for %s (%s) to %s\n", r.Id, t.Id, c.card.Issuer) // change to logs
		fmt.Println("Verifying transaction...")
		fmt.Println("Transaction approved, processing payment...")

		t.Status = "pending"
		t.ConfirmationCode = NewUUID()

		err = c.card.Credit(r.Amount)

		if err != nil {
			log.Fatal(err) // handle this better
		} else {
			fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
			fmt.Printf("Refund request successful you have been credited %f\n", r.Amount)

			t.Status = "paid"
			r.Status = "cancelled"

			err = Save(r)
		}
	}

}

func (c *CreditCardProcessor) Invoice(t *Transaction) {
	issueInvoice(t, c.Label)
}

func issueInvoice(t *Transaction, p string) {

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
	fmt.Printf(i, p, t.User.FullName, t.Id, NewUUID(), t.Amount, t.Amount,
		date, t.Status, t.ConfirmationCode)
}
