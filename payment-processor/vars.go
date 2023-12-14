package main

import "time"

// files

var UserFile File
var TransactionFile File
var CardFile File

// billing period

const (
	week  = 168 * time.Hour
	month = 672 * time.Hour
	year  = 8064 * time.Hour
)

// paypal

const (
	// product
	createProduct = "https://api-m.sandbox.paypal.com/v1/catalogs/products"

	// subs
	createPlan = "https://api-m.sandbox.paypal.com/v1/billing/plans"
	getPlan    = "https://api-m.sandbox.paypal.com/v1/billing/plans/"
	CreateSub  = "https://api-m.sandbox.paypal.com/v1/billing/subscriptions"
	GetSub     = "https://api-m.sandbox.paypal.com/v1/billing/subscriptions/"

	// auth
	getToken = ""

	// payments
	create  = "https://paypal/api/v1/payment/create"
	approve = "https://paypal/api/v1/payment/approve"
	execute = "https://paypal/api/v1/payment/execute"
)
