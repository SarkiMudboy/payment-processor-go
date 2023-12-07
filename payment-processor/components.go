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
	Tag() string
}

func (u *User) String() string {
	return fmt.Sprintf("%s (%s)", u.FullName, u.Id)
}

func (u *User) Get() string {
	return u.Id
}

func (u *User) Tag() string {
	return "users"
}

func (t *Transaction) String() string {
	return fmt.Sprintf("Transaction: %s (%f)", t.Id, t.Amount)
}

func (t *Transaction) Get() string {
	return t.Id
}

func (t *Transaction) Tag() string {
	return "transactions"
}

func Save(c Component) error {
	// json loading from Db

	db := make(map[string]interface{})

	u, _ := UserFile.Open()
	defer u.Close()

	entries, err := ioutil.ReadAll(u)
	fmt.Println(entries)
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

	_, err = u.Write([]byte(jsonData))

	if err != nil {
		return fmt.Errorf("Cannot save data: %s", err)
	}

	_ = u.Close()

	return nil
}

func Load(filename string) (map[string]Component, error) {
	// json reading from db

	var db map[string]Component

	u, _ := UserFile.Open()
	defer u.Close()

	entries, err := ioutil.ReadAll(u)

	if err != nil {
		return nil, fmt.Errorf("Could not open file: %s\n", err)
	}

	if len(entries) != 0 {

		err = json.Unmarshal(entries, &db)

		if err != nil {
			return nil, fmt.Errorf("Unable to load data %s\n", err)
		}

	}

	// return the data
	return db, nil
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
