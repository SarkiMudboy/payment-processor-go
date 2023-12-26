module payments-go/payment-processor/main

go 1.21.1

replace payment-go/payment-processor/payments => ../payments

require payment-go/payment-processor/payments v0.0.0-00010101000000-000000000000

require (
	github.com/google/uuid v1.5.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
)
