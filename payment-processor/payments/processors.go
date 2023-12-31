package payments

import (
	"fmt"
	"log"
)

type Processor interface {
	OneTimePayment(*Transaction)
	Subscription(*subscription, *Transaction)
	Refund(*Transaction, *Transaction)
	Invoice(*Transaction)
}

type CreditCardProcessor struct {
	card
	Label string
}

func (c *CreditCardProcessor) OneTimePayment(t *Transaction) {

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

func (c *CreditCardProcessor) Subscription(s *subscription, t *Transaction) {

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
		log.Println(failedTransactionError)
	}

	c.Invoice(t)

}

func (c *CreditCardProcessor) Refund(r, t *Transaction) {

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
			log.Println(failedTransactionError, err)
			t.Status = "failed"
		} else {
			fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
			fmt.Printf("Refund request successful you have been credited %f\n", r.Amount)

			t.Status = "paid"
			r.Status = "cancelled"

			err = Save(&TransactionFile, r)
		}

		c.Invoice(t)
	}

}

func (c *CreditCardProcessor) Invoice(t *Transaction) {
	issueInvoice(t, c.Label)
}

// account

type BankAccountProcessor struct {
	account
	Label string
}

func (b *BankAccountProcessor) OneTimePayment(t *Transaction) {

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

func (b *BankAccountProcessor) Subscription(s *subscription, t *Transaction) {

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

func (b *BankAccountProcessor) Refund(r, t *Transaction) {

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

			err = Save(&TransactionFile, r)
		}

		b.Invoice(t)
	}
}

func (b *BankAccountProcessor) Invoice(t *Transaction) {
	issueInvoice(t, b.Label)
}

// PayPal

type PayPalProcessor struct {
	client paypalClient
	Label  string
}

func (p *PayPalProcessor) OneTimePayment(t *paypalTransaction) {

	err := p.Pay(t)

	if err != nil {
		t.Status = "failed"
		fmt.Printf("Transaction failed, contact paypal support. Transaction ID: %s\n", t.Id)
	}

	t.ConfirmationCode = NewUUID()
	fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)

}

func (p *PayPalProcessor) Subscription(s paypalSubscription, t *paypalTransaction) {

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

func (p PayPalProcessor) Invoice(t *paypalTransaction) {

	transaction_data := make(map[string]interface{}, 3)
	ToMap(t, transaction_data)
	t.genrateInvoice(transaction_data)

}

// paypal refund

func (p PayPalProcessor) Refund(r, t *paypalTransaction) {

	entries, err := Load(TransactionFile)

	if err != nil {
		log.Fatal(err)
	}

	_, err = r.Load(entries)

	if err != nil {
		log.Fatal("Could not find transaction:", err)
	} else {
		fmt.Printf("submitting refund request for %s (%s) to paypal\n", r.Id, t.Id) // change to logs
		fmt.Println("Verifying transaction...")
		fmt.Println("Transaction approved, processing payment...")

		t.Status = "pending"
		t.ConfirmationCode = NewUUID()

		refundData := []string{sellerID, p.client.creds[0]}

		err := t.requestRefund(p.client, refundData)

		if err != nil {
			log.Fatal(err) // handle this better
			t.Status = "failed"
		} else {
			fmt.Printf("(Success) Transaction complete, your confirmation code is %s\n", t.ConfirmationCode)
			fmt.Printf("Refund request successful you have been credited %f\n", r.Amount)

			t.Status = "paid"
			r.Status = "cancelled"

			err = Save(&TransactionFile, r)
		}

		p.Invoice(t)
	}
}
