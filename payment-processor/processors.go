package main

import (
	"errors"
	"fmt"
	"log"
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
	card
	Label string
}

func (c *CreditCardProcessor) oneTimePayment(t *transaction) {

	fmt.Printf("submitting transaction: %s to %s\n", t.Id, c.Issuer) // change to logs
	fmt.Println("Verifying transaction...")
	fmt.Println("Transaction approved, processing payment...")

	t.Status = "pending"
	t.ConfirmationCode = NewUUID()

	err := c.Charge(t.Amount)

	if err != nil {
		t.Status = "failed"
		log.Fatal("(Error) Could not process payment:", err)
	} else {
		fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
		t.Status = "paid"
	}

	c.Invoice(t)

}

func (c *CreditCardProcessor) Subscription(s *subscription, t *transaction) {

	if s.VerifyBilling() {
		fmt.Printf("submitting transaction: %s to %s\n", t.Id, c.Issuer) // change to logs
		fmt.Println("Verifying transaction...")
		fmt.Println("Transaction approved, processing payment...")

		t.Status = "pending"
		t.ConfirmationCode = NewUUID()

		err := c.Charge(s.plan)

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

func (c *CreditCardProcessor) Refund(r, t *transaction) {

	entries, err := Load(TransactionFile)

	if err != nil {
		log.Fatal(err)
	}

	_, err = r.Load(entries)

	if err != nil {
		log.Fatal("Could not find transaction:", err)
	} else {
		fmt.Printf("submitting refund request for %s (%s) to %s\n", r.Id, t.Id, c.Issuer) // change to logs
		fmt.Println("Verifying transaction...")
		fmt.Println("Transaction approved, processing payment...")

		t.Status = "pending"
		t.ConfirmationCode = NewUUID()

		err = c.Credit(r.Amount)

		if err != nil {
			log.Fatal(err) // handle this better
			t.Status = "failed"
		} else {
			fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
			fmt.Printf("Refund request successful you have been credited %f\n", r.Amount)

			t.Status = "paid"
			r.Status = "cancelled"

			err = Save(r)
		}

		c.Invoice(t)
	}

}

func (c *CreditCardProcessor) Invoice(t *transaction) {
	issueInvoice(t, c.Label)
}

type account struct {
	User    user    `json:"user"`
	Holder  string  `json:"holder"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
	Bank    string  `json:"bank"`
}

func (a *account) Debit(amount float64) error {
	if a.Balance > amount {
		a.Balance -= amount

		return nil
	}

	return errors.New("insufficient funds!")
}

func (a *account) Credit(amount float64) error {
	a.Balance += amount

	return nil
}

type BankAccountProcessor struct {
	account
	Label string
}

func (b *BankAccountProcessor) oneTimePayment(t *transaction) {

	fmt.Printf("submitting transaction: %s to %s\n", t.Id, b.Bank) // change to logs
	fmt.Println("Verifying transaction...")
	fmt.Println("Transaction approved, processing payment...")

	t.Status = "pending"
	t.ConfirmationCode = NewUUID()

	err := b.Debit(t.Amount)

	if err != nil {
		t.Status = "failed"
		log.Fatal("(Error) Could not process payment: ", err)
	} else {
		fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
		t.Status = "paid"
	}

	b.Invoice(t)
}

func (b *BankAccountProcessor) Subscription(t *transaction, s *subscription) {

	if s.VerifyBilling() {
		fmt.Printf("submitting transaction: %s to %s\n", t.Id, b.Bank) // change to logs
		fmt.Println("Verifying transaction...")
		fmt.Println("Transaction approved, processing payment...")

		t.Status = "pending"
		t.ConfirmationCode = NewUUID()

		err := b.Debit(s.plan)

		if err != nil {
			t.Status = "failed"
			log.Fatal("(Error) Could not process payment: ", err)
		} else {
			fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
			fmt.Printf("You have renewed your %s subscription plan till %s\n", s.name, s.due)
			t.Status = "paid"
		}
	} else {
		fmt.Println("(Fail) Transaction failed")
	}

	b.Invoice(t)
}

func (b *BankAccountProcessor) Refund(r, t *transaction) {

	entries, err := Load(TransactionFile)

	if err != nil {
		log.Fatal(err)
	}

	_, err = r.Load(entries)

	if err != nil {
		log.Fatal("Could not find transaction:", err)
	} else {
		fmt.Printf("submitting refund request for %s (%s) to %s\n", r.Id, t.Id, b.Bank) // change to logs
		fmt.Println("Verifying transaction...")
		fmt.Println("Transaction approved, processing payment...")

		t.Status = "pending"
		t.ConfirmationCode = NewUUID()

		err := b.Credit(r.Amount)

		if err != nil {
			log.Fatal(err) // handle this better
			t.Status = "failed"
		} else {
			fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
			fmt.Printf("Refund request successful you have been credited %f\n", r.Amount)

			t.Status = "paid"
			r.Status = "cancelled"

			err = Save(r)
		}

		b.Invoice(t)
	}
}

func (b *BankAccountProcessor) Invoice(t *transaction) {
	issueInvoice(t, b.Label)
}

// PayPal

type PayPalProcessor struct {
	client paypalClient
	Label  string
}

func (p *PayPalProcessor) oneTimePayment(t *paypalTransaction) {

	err := p.Pay(t)

	if err != nil {
		t.Status = "failed"
		fmt.Printf("Transaction failed, contact paypal support. Transaction ID: %s\n", t.Id)
	}

	t.ConfirmationCode = NewUUID()
	fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)

}

func (p *PayPalProcessor) subscription(s paypalSubscription, t *paypalTransaction) {

	if s.VerifyBilling(p.client) {
		fmt.Printf("submitting transaction: %s to paypal \n", t.Id) // change to logs

		t.Status = "pending"
		t.ConfirmationCode = NewUUID()

		err := p.Pay(t)

		if err != nil {
			t.Status = "failed"
			log.Fatal("(Error) Could not process payment: ", err)
		} else {
			fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
			fmt.Printf("You have renewed your %s subscription plan till %s\n", s.name, s.due)
			t.Status = "paid"
		}
	} else {
		fmt.Println("(Fail) Transaction failed")
	}

	p.Invoice(t)
}

func (p *PayPalProcessor) Pay(t *paypalTransaction) error {

	data := make(map[string]interface{}, 3)
	data = ToMap(t, data)

	fmt.Println("Creating order...")

	// create the order
	t.createOrder(data)

	fmt.Println("Order created, confirming...")
	t.Status = "pending"

	// confirm the order
	confirmed := t.confirmOrder(data)
	if !confirmed {
		return fmt.Errorf("Cannot confirm paypal order")
	}

	fmt.Println("Confirmed, Authorizing payment....")

	// authorize the order
	authorized := t.authorizeOrder()
	if !authorized {
		return fmt.Errorf("from paypal: Unauthorized!")
	}
	t.Status = "paid"
	return nil
}

func (p *PayPalProcessor) Invoice(t *paypalTransaction) {

}
