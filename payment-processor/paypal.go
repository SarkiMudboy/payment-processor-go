package main

import (
	"encoding/json"
	"log"
)

type paypalTransaction struct {
	transaction
	order  string
	client paypalClient
}

func (p *paypalTransaction) createOrder(data map[string]interface{}) {

}

type paypalSubscription struct {
	subscription
	subId   string
	planId  string
	product map[string]interface{}
	client  paypalClient
}

type paypalClient struct {
	user  user
	token string
	creds []string
}

func (p paypalClient) Headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + p.token,
		"Content-Type":  "application/json",
	}
}

func (p paypalClient) GetToken() {

}

func (p *paypalSubscription) getProduct() map[string]interface{} {

	b, err := json.Marshal(p.product)

	if err != nil {
		log.Fatal(err)
	}

	response := Request(createProduct, "GET", b, p.client.Headers())

	if response["status"] == 200 {
		return response["plan"].(map[string]interface{})
	}

	return nil
}

func (p *paypalSubscription) CreatePlan(d map[string]interface{}) {

	b, err := json.Marshal(d)

	if err != nil {
		log.Fatal(err)
	}

	response := Request(createPlan, "POST", b, p.client.Headers())

	if response["status"] != 201 {
		log.Fatal("Error creating plan")
	}

	plan, ok := response["plan"].(map[string]interface{})

	if !ok {
		log.Fatal("Invalid data")
	}

	plan_id, ok := plan["id"].(string)

	if !ok {
		log.Fatal("Invalid data")
	}

	p.planId = plan_id

}

func (p *paypalSubscription) GetPlan() map[string]interface{} {

	endpoint := getPlan + p.planId

	response := Request(endpoint, "POST", []byte{}, p.client.Headers())

	if response["status"] == 200 {
		return response["plan"].(map[string]interface{})
	}

	return nil
}

func (p *paypalSubscription) CreateSub(data map[string]interface{}) {

	b, err := json.Marshal(&data)

	if err != nil {
		log.Fatal(err)
	}

	response := Request(CreateSub, "POST", b, p.client.Headers())

	if response["status"] != 201 {
		log.Fatal("Error creating sub")
	}

	sub, ok := response["subscription"].(map[string]interface{})

	if !ok {
		log.Fatal("Invalid data")
	}

	id, ok := sub["id"].(string)

	if !ok {
		log.Fatal("Invalid data!")
	}

	p.subId = id
}

func (p *paypalSubscription) GetSub() map[string]interface{} {

	endpoint := GetSub + p.subId

	response := Request(endpoint, "POST", []byte{}, p.client.Headers())

	if response["status"] == 200 {
		return response["subscription"].(map[string]interface{})
	}

	return nil

}

func (p *paypalSubscription) VerifyBilling(c paypalClient) {

}

func (p *paypalSubscription) Bill(c paypalClient) {

}
