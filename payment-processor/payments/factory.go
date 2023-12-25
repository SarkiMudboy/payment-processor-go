package payments

import (
	"log"
	"os"
	"time"
)

func InitUser(id string) user {
	return user{Id: id}
}
func InitTransaction(id string) transaction {
	return transaction{Id: id}
}
func InitCard(id string) card {
	return card{Id: id}
}
func InitSubscription(id string) subscription {
	return subscription{id: id}
}

func NewUser(id, username, firstname, lastname string) user {

	u := user{
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
	}

	u.Id = NewUUID()
	u.CreatedAt = time.Now()

	return u
}

func NewTransaction(user user, amount float64, status string) transaction {

	t := transaction{
		User:   user,
		Amount: amount,
		Status: status,
	}

	t.Id = NewUUID()
	t.CreatedAt = time.Now()

	return t
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

	_, err := os.OpenFile("db/"+filename+"."+filetype, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		log.Fatal(err)
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
