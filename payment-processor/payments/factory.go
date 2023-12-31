package payments

import (
	"fmt"
	"log"
	"os"
	"time"
)

func InitUser(id string) user {
	return user{Id: id}
}
func InitTransaction(id string) Transaction {
	return Transaction{Id: id}
}
func InitCard(id string) card {
	return card{Id: id}
}
func InitSubscription(id string) subscription {
	return subscription{id: id}
}

func NewCreditCardProcessor(creditCard card, label string) *CreditCardProcessor {
	return &CreditCardProcessor{
		card:  creditCard,
		Label: label,
	}
}

func NewBankAccountProcessor(a account, label string) *BankAccountProcessor {
	return &BankAccountProcessor{
		account: a,
		Label:   label,
	}
}

func NewPayPalProcessor(c paypalClient, label string) *PayPalProcessor {
	return &PayPalProcessor{
		client: c,
		Label:  label,
	}
}

func NewUser(username, firstname, lastname string) user {

	u := user{
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
	}

	u.Id = NewUUID()
	u.FullName = u.FirstName + " " + u.LastName
	u.CreatedAt = time.Now()

	return u
}

func NewTransaction(user user, amount float64, status string) *Transaction {

	t := &Transaction{
		User:   user,
		Amount: amount,
		Status: status,
	}

	t.Id = NewUUID()
	t.CreatedAt = time.Now()

	return t
}

func NewAccount(user user, holder string, number string, bank string) *account {
	a := &account{
		User:   user,
		Holder: holder,
		Number: number,
		Bank:   bank,
	}

	a.Id = NewUUID()
	return a
}

func NewCard(user user, issuer string, number string, expiry string, cvv string) card {

	expiry_parsed, err := time.Parse("2006-01-02 03:04:05", expiry)

	if err != nil {
		log.Fatal(err)
	}

	c := card{
		User:   user,
		Issuer: issuer,
		Number: number,
		Expiry: expiry_parsed,
		CVV:    cvv,
	}

	return c
}

func NewFile(filename, filetype string) (File, error) {

	f := File{
		Name: filename,
		Type: filetype,
	}

	filepath := filename + "." + filetype

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		file, err := os.Create(filepath)
		fmt.Println(file)

		defer file.Close()

		if err != nil {
			return File{}, err
		}

	} else {
		file, err := os.Open(filepath)
		defer file.Close()

		if err != nil {
			return File{}, err
		}
	}

	return f, nil
}

func NewSubscription(user user, name string, period string, plan float64) subscription {

	s := subscription{
		user:   user,
		name:   name,
		period: period,
		plan:   plan,
	}

	s.active = true
	s.SetBilling()

	return s
}

func NewPaypalSub(sub subscription) paypalSubscription {
	p := paypalSubscription{
		subscription: sub,
	}

	// some stuff to do here: ...
	return p
}
