package payments

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
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

func ToMap(from interface{}, to map[string]interface{}) map[string]interface{} {

	err := mapstructure.Decode(from, &to)

	if err != nil {
		log.Fatal(err)
	}

	return to
}
