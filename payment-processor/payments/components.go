package payments

import (
	"encoding/json"
	"fmt"
	"io"
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
	PendingTransactions []Transaction `json:"pending transactions,omitempty"`
	TotalPaidAmount     float64       `json:"total paid amount,omitempty"`
	CreatedAt           time.Time     `json:"time"`
}

type Transaction struct {
	Id               string    `json:"id"`
	ConfirmationCode string    `json:"conf_code,omitempty"`
	User             user      `json:"user"`
	Amount           float64   `json:"amount,omitempty"`
	Status           string    `json:"status"`
	Invoice          string    `json:"invoice,omitempty"`
	CreatedAt        time.Time `json:"time"`
}

func init() {
	UserFile, _ = NewFile("users", "json")
	TransactionFile, _ = NewFile("transactions", "json")
	AccountFile, _ = NewFile("accounts", "json")
}

type Component interface {
	String() string
	Get() string
}

func (u user) String() string {
	return fmt.Sprintf("%s (%s)", u.FullName, u.Id)
}

func (u user) Get() string {
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

func (t Transaction) String() string {
	return fmt.Sprintf("Transaction: %s (%f)", t.Id, t.Amount)
}

func (t Transaction) Get() string {
	return t.Id
}

func (t *Transaction) Load(b []byte) (Transaction, error) {

	var db map[string]Transaction

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

func Save(file *File, c Component) error {
	// saving updated data (json) to Db

	db := make(map[string]interface{})

	f, _ := file.Open()
	defer f.Close()

	entries, err := io.ReadAll(f)

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

	err = f.Truncate(0)
	_, err = f.Seek(0, 0)

	_, err = f.Write([]byte(jsonData))

	if err != nil {
		return fmt.Errorf("Cannot save data: %s", err)
	}

	_ = f.Close()

	return nil
}

func Load(file File) ([]byte, error) {
	// json reading from db as bytes (pair with component.Load() function to get data)

	u, _ := file.Open()
	defer u.Close()

	data, err := io.ReadAll(u)

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

type account struct {
	Id      string  `json:"Id"`
	User    user    `json:"user"`
	Holder  string  `json:"holder"`
	Number  string  `json:"number"`
	Balance float64 `json:"balance"`
	Bank    string  `json:"bank"`
}

func (a account) String() string {
	return fmt.Sprintf("%s (%s)", a.User.FullName, a.Number)
}

func (a account) Get() string {
	return a.Id
}

func (a *account) Load(b []byte) (account, error) {

	var db map[string]account

	err := json.Unmarshal(b, &db)

	if err != nil {
		log.Fatal(err)
	}

	account, ok := db[a.Id]

	if !ok {
		return account, fmt.Errorf("account does not exist!")
	}

	return account, nil
}

func (a *account) Debit(amount float64) error {
	if a.Balance > amount {
		a.Balance -= amount

		return nil
	}

	return insufficientError
}

func (a *account) Credit(amount float64) error {
	a.Balance += amount

	return nil
}

func (a *account) transfer() {} //decide if to put in processor?

type bankTransaction struct {
	Id              string  `json:"id"`
	Sender          user    `json:"sender"`
	Reciepient      string  `json:"reciepient"` //user id string
	Receipt         string  `json:"receipt"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"type"`
}

func (b bankTransaction) Get() string {
	return b.Id
}

func (b bankTransaction) String() string {
	return fmt.Sprintf("%s (%f)", b.Id, b.Amount)
}

func (t *bankTransaction) Load(b []byte) (bankTransaction, error) {

	var db map[string]bankTransaction

	err := json.Unmarshal(b, &db)

	if err != nil {
		log.Fatal(err)
	}

	transaction, ok := db[t.Id]

	if !ok {
		return transaction, fmt.Errorf("transaction does not exist!")
	}

	return transaction, nil
}

func (t *bankTransaction) issueReceipt() {

	date := time.Now().Format(time.RFC3339)

	r := `
	-----------(%s)---------------
	Name: %s
	Transaction: %s
	Receipt Number: %s
	------------------------------
	Amount: %f
	Date: %s
	`
	r = fmt.Sprintf(r, t.Sender.FullName, t.Id, t.Amount,
		date)

	fmt.Println(r)

	t.Receipt = r
}

type File struct {
	Name string
	Type string
	Dir  string
}

func (f *File) Open() (*os.File, error) {
	// fmt.Println(f)

	file, err := os.OpenFile(f.Dir, os.O_RDWR|os.O_CREATE, 0755)

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
