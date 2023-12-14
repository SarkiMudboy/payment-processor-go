package main

import (
	"fmt"

	"github.com/google/uuid"
)

func NewUUID() string {

	id := uuid.New()
	return id.String()

}

func Request(endpoint string, method string, body []byte, headers map[string]string) map[string]interface{} {
	fmt.Printf("[%s] pinging %s .....\n", method, endpoint)
	fmt.Println("OK server response")

	return make(map[string]interface{})
}
