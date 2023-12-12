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
	createPlan = "http://paypal/api/v1/payment/subscription/create"
	create     = "https://paypal/api/v1/payment/create"
	approve    = "https://paypal/api/v1/payment/approve"
	execute    = "https://paypal/api/v1/payment/execute"
)
