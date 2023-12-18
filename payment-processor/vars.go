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
	getToken = "https://api-m.sandbox.paypal.com/v1/oauth2/token"

	// payments
	createOrder    = "https://api-m.sandbox.paypal.com/v2/checkout/orders"
	confirmOrder   = `https://api-m.sandbox.paypal.com/v2/checkout/orders/%s/confirm-payment-source`
	authorizeOrder = `https://api-m.sandbox.paypal.com/v2/checkout/orders/%s/authorize`
)
