package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type User struct {
	Id                  string        `json:"id"`
	Username            string        `json:"username"`
	FirstName           string        `json:"first name"`
	LastName            string        `json:"last name"`
	FullName            string        `json:"fullname,omitempty"`
	PendingTransactions []Transaction `json:"pending transactions,omitempty"`
	TotalPaidAmount     float64       `json:"total paid amount,omitempty"`
	CreatedAt           time.Time     `json:"time"`
}

type Transaction struct {
	Id               string    `json:"id"`
	ConfirmationCode string    `json:"conf_code,omitempty"`
	User             User      `json:"user"`
	Amount           float64   `json:"amount,omitempty"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"time"`
}

func init() {
	UserFile, _ = NewFile("users", "json")
	TransactionFile, _ = NewFile("transactions", "json")
}

func NewUser(username, firstname, lastname string) *User {

	u := &User{
		Username:  username,
		FirstName: firstname,
		LastName:  lastname,
	}

	u.Id = NewUUID()
	u.CreatedAt = time.Now()

	return u
}

func NewTransaction(user User, amount float64, status string) *Transaction {

	t := &Transaction{
		User:   user,
		Amount: amount,
		Status: status,
	}

	t.Id = NewUUID()
	t.CreatedAt = time.Now()

	return t
}

type Component interface {
	String() string
	Get() string
	Load([]byte) (Component, error)
}

func (u *User) String() string {
	return fmt.Sprintf("%s (%s)", u.FullName, u.Id)
}

func (u *User) Get() string {
	return u.Id
}

func (u *User) Load(b []byte) (*User, error) {

	var db map[string]*User

	err := json.Unmarshal(b, &db)

	if err != nil {
		log.Fatal(err)
	}

	user, ok := db[u.Id]

	if !ok {
		return nil, fmt.Errorf("User does not exist!")
	}

	return user, nil
}

func (t *Transaction) String() string {
	return fmt.Sprintf("Transaction: %s (%f)", t.Id, t.Amount)
}

func (t *Transaction) Get() string {
	return t.Id
}

func (t *Transaction) Load(b []byte) (*Transaction, error) {

	var db map[string]*Transaction

	err := json.Unmarshal(b, &db)

	if err != nil {
		log.Fatal(err)
	}

	transaction, ok := db[t.Id]

	if !ok {
		return nil, fmt.Errorf("Transaction does not exist!")
	}

	return transaction, nil
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

func NewFile(filename, filetype string) (File, error) {

	f := File{
		Name: filename,
		Type: filetype,
	}

	_, err := os.OpenFile(filename+"."+filetype, os.O_RDWR|os.O_CREATE, 0755)

	if err != nil {
		log.Fatal(err)
	}

	return f, nil
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
