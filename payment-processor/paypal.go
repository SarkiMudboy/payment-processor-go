package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type paypalTransaction struct {
	transaction
	order     string
	client    paypalClient
	captureID string
}

func (p *paypalTransaction) createOrder(data map[string]interface{}) {

	b, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	response := Request(createOrder, "POST", b, p.client.Headers())

	if response["status"] != 201 {
		log.Fatal(getResourceError)
	}

	order, ok := response["plan"].(map[string]interface{})

	if !ok {
		log.Fatal(invalidResource)
	}

	order_id, ok := order["id"].(string)

	if !ok {
		log.Fatal(invalidResource)
	}

	p.order = order_id

}

func (p *paypalTransaction) confirmOrder(data map[string]interface{}) bool {

	b, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	endpoint := fmt.Sprintf(confirmOrder, p.order)

	response := Request(endpoint, "POST", b, p.client.Headers())

	if response["status"] != 200 {
		log.Println(getResourceError)
		return false
	}

	return true
}

func (p *paypalTransaction) authorizeOrder() bool {

	endpoint := fmt.Sprintf(authorizeOrder, p.order)

	response := Request(endpoint, "POST", []byte{}, p.client.Headers())

	if response["status"] != 200 {
		log.Println(failedError)
		return false
	}

	return true

}

func (p *paypalTransaction) genrateInvoice(data map[string]interface{}) {
	// generate and send invoice

	b, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	response := Request(createInvoice, "POST", b, p.client.Headers())

	if response["status"] != 201 {
		log.Fatal(createResourceError)
	}

	invoice, ok := response["invoice"].(map[string]interface{})

	if !ok {
		log.Fatal(invalidResource)
	}

	invoice_id, ok := invoice["id"].(string)
	if !ok {
		log.Fatal(invalidResource)
	}

	p.Invoice = invoice_id

	// send
	endpoint := fmt.Sprintf(sendInvoice, p.Invoice)
	response = Request(endpoint, "POST", b, p.client.Headers())

	if response["status"] != 200 {
		log.Fatal(failedError)
	}

}

func (p paypalTransaction) requestRefund(c paypalClient, refundData []string) error {
	buildHeader := c.getAssertionValue(refundData...)

	headers := p.client.Headers()
	headers["PayPal-Auth-Assertion"] = buildHeader

	endpoint := fmt.Sprintf(refund, p.captureID)
	response := Request(endpoint, "POST", []byte{}, headers)

	if response["status"] != 200 {
		return failedError
	}

	return nil
}

type paypalSubscription struct {
	subscription
	subId   string
	planId  string
	product map[string]interface{}
	client  paypalClient
}

type paypalClient struct {
	user          user
	token         string
	creds         []string
	authAssertion string
}

func (p *paypalClient) getAssertionValue(data ...string) string {
	client, seller := data[0], data[1]

	val := fmt.Sprintf("%s/%s", client, seller)
	return val

}

func (p paypalClient) Headers() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + p.token,
		"Content-Type":  "application/json",
	}
}

func (p *paypalClient) GetToken() error {
	clientId, clientSecret := p.creds[0], p.creds[1]

	data := map[string]string{
		"CLIENT_ID":     clientId,
		"CLIENT_SECRET": clientSecret,
	}

	body, err := json.Marshal(data)

	if err != nil {
		return invalidResource
	}

	response := Request(getToken, "POST", body, p.Headers())

	if response["status"] != 200 {
		return getResourceError
	}

	token, ok := response["token"].(string)

	if !ok {
		return invalidResource
	}

	p.token = token

	return nil
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
		log.Fatal(createResourceError)
	}

	plan, ok := response["plan"].(map[string]interface{})

	if !ok {
		log.Fatal(invalidResource)
	}

	plan_id, ok := plan["id"].(string)

	if !ok {
		log.Fatal(invalidResource)
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

	b, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	response := Request(CreateSub, "POST", b, p.client.Headers())

	if response["status"] != 201 {
		log.Fatal(createResourceError)
	}

	sub, ok := response["subscription"].(map[string]interface{})

	if !ok {
		log.Fatal(invalidResource)
	}

	id, ok := sub["id"].(string)

	if !ok {
		log.Fatal(invalidResource)
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

func (p *paypalSubscription) VerifyBilling(c paypalClient) bool {
	return true
}

func (p *paypalSubscription) Bill(c paypalClient) bool {
	return true
}
