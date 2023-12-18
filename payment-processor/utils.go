package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

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
	fmt.Printf(i, p, t.User.FullName, t.Id, NewUUID(), t.Amount, t.Amount,
		date, t.Status, t.ConfirmationCode)
}

func NewUUID() string {

	id := uuid.New()
	return id.String()

}

func Request(endpoint string, method string, body []byte, headers map[string]string) map[string]interface{} {
	fmt.Printf("[%s] pinging %s .....\n", method, endpoint)
	fmt.Println("OK server response")

	return make(map[string]interface{})
}

func ToMap(from interface{}, to map[string]interface{}) map[string]interface{} {

	err := mapstructure.Decode(from, &to)

	if err != nil {
		log.Fatal(err)
	}

	return to
}
