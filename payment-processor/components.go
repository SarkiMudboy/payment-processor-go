package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type user struct {
	Id                  string        `json:"id"`
	Username            string        `json:"username"`
	FirstName           string        `json:"first name"`
	LastName            string        `json:"last name"`
	FullName            string        `json:"fullname,omitempty"`
	PendingTransactions []transaction `json:"pending transactions,omitempty"`
	TotalPaidAmount     float64       `json:"total paid amount,omitempty"`
	CreatedAt           time.Time     `json:"time"`
}

type transaction struct {
	Id               string    `json:"id"`
	ConfirmationCode string    `json:"conf_code,omitempty"`
	User             user      `json:"user"`
	Amount           float64   `json:"amount,omitempty"`
	Status           string    `json:"status"`
	Invoice          string    `json:"invoice,omitempty"`
	CreatedAt        time.Time `json:"time"`
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

func init() {
	UserFile, _ = NewFile("users", "json")
	TransactionFile, _ = NewFile("transactions", "json")
	CardFile, _ = NewFile("cards", "json")
}

type Component interface {
	String() string
	Get() string
}

func (u *user) String() string {
	return fmt.Sprintf("%s (%s)", u.FullName, u.Id)
}

func (u *user) Get() string {
	return u.Id
}

func (u *user) Load(b []byte) (user, error) {

	var db map[string]user

	err := json.Unmarshal(b, &db)

	if err != nil {
		log.Fatal(err)
	}

	user, ok := db[u.Id]

	if !ok {
		return user, fmt.Errorf("User does not exist!")
	}

	return user, nil
}

func (t *transaction) String() string {
	return fmt.Sprintf("Transaction: %s (%f)", t.Id, t.Amount)
}

func (t *transaction) Get() string {
	return t.Id
}

func (t *transaction) Load(b []byte) (transaction, error) {

	var db map[string]transaction

	err := json.Unmarshal(b, &db)

	if err != nil {
		log.Fatal(err)
	}

	transaction, ok := db[t.Id]

	if !ok {
		return transaction, fmt.Errorf("Transaction does not exist!")
	}

	return transaction, nil
}

func (c *card) String() string {
	return fmt.Sprintf("%s (%s)", c.User.FullName, c.Issuer)
}

func (c *card) Get() string {
	return c.Id
}

func (c *card) Load(b []byte) (card, error) {

	var db map[string]card

	err := json.Unmarshal(b, &db)

	if err != nil {
		log.Fatal(err)
	}

	card, ok := db[c.Id]

	if !ok {
		return card, fmt.Errorf("Transaction does not exist!")
	}

	return card, nil
}

func (c *card) Charge(amount float64) error {

	if !c.Expired() {
		if c.Balance > amount && amount < c.Limit {
			c.Balance -= amount
			return nil
		}

		return errors.New("Insufficient funds!")
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

func Save(c Component) error {
	// saving updated data (json) to Db

	db := make(map[string]interface{})

	u, _ := UserFile.Open()
	defer u.Close()

	entries, err := ioutil.ReadAll(u)

	if err != nil {
		return fmt.Errorf("Could not open file: %s\n", err)
	}

	if len(entries) != 0 {

		err = json.Unmarshal(entries, &db)

		if err != nil {
			return fmt.Errorf("Unable to load data %s\n", err)
		}
	}

	db[c.Get()] = c
	fmt.Println("Updated map data", db)

	jsonData, err := json.MarshalIndent(db, "", " ")

	if err != nil {
		return fmt.Errorf("Cannot marshal (serialize) data: %s", err)
	}

	err = u.Truncate(0)
	_, err = u.Seek(0, 0)

	_, err = u.Write([]byte(jsonData))

	if err != nil {
		return fmt.Errorf("Cannot save data: %s", err)
	}

	_ = u.Close()

	return nil
}

func Load(file File) ([]byte, error) {
	// json reading from db as bytes (pair with component.Load() function to get data)

	u, _ := file.Open()
	defer u.Close()

	data, err := ioutil.ReadAll(u)

	if err != nil {
		return nil, fmt.Errorf("Could not open file: %s\n", err)
	}

	if len(data) != 0 {
		return data, nil
	}

	err = u.Close()

	// return error if file is empty
	return nil, fmt.Errorf("Empty file: %s\n", err)
}

type File struct {
	Name string
	Type string
}

func (f *File) Open() (*os.File, error) {

	file, err := os.OpenFile(f.Name+"."+f.Type, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		return nil, fmt.Errorf("%s\n", err)
	}

	return file, nil
}

func (f *File) Close(file *os.File) error {

	err := file.Close()

	if err != nil {
		return fmt.Errorf("%s\n", err)
	}

	return nil

}
