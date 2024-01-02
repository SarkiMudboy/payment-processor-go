package main

// use flagset from flag package
// go-wallet (file name) + pay/send/mod (subcommand) + -op/-ref/-sub (flagset) + -amount (flagset) + $600 (float64 var)
// > card
// > bank
// (Auth) user ID -> ""

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"payment-go/payment-processor/payments"
)

var card = map[string]string{
	"issuer": "Verve",
	"number": "123-4567-8901-234-5678",
	"expiry": "2025-08-15T14:30:45.0000001-05:00",
	"cvv":    "987",
}

var account = map[string]string{
	"holder": "ME",
	"number": "3121526707",
	"bank":   "First Bank",
}

var proccessorOptions = `
	Pay via?
	 > card (1)
	 > bank (2)

	 Enter (num) to continue V:
`

func NewPayCommand() *PayCommand {
	pay := &PayCommand{
		fs: flag.NewFlagSet("pay", flag.ContinueOnError),
	}

	pay.fs.StringVar(&pay.operation, "operation", "op", "represents the payment operation")
	pay.fs.Float64Var(&pay.amount, "amount", 0.0, "transaction amount")

	return pay
}

type PayCommand struct {
	fs        *flag.FlagSet
	operation string
	amount    float64
}

func (p *PayCommand) Name() string {
	return p.fs.Name()
}

func (p *PayCommand) Init(args []string) error {
	return p.fs.Parse(args)
}

func (p *PayCommand) Run() error {
	req := GetRequest(p)

	handler := req.Build()
	return handler.Handle(req)
}

func (p *PayCommand) Operation() string {
	return p.operation
}

type Runner interface {
	Init([]string) error
	Run() error
	Name() string
	Operation() string
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("Enter a subcommand")
	}

	cmds := []Runner{
		NewPayCommand(),
	}

	subcommand := os.Args[1]

	for _, command := range cmds {
		if command.Name() == subcommand {
			command.Init(os.Args[2:])
			return command.Run()
		}
	}

	return fmt.Errorf("Unknown Command: %s", subcommand)
}

func GetRequest(r Runner) Request {
	if r.Name() == "pay" {

		var userID string
		var proc int
		var processor string
		var data struct {
			card         map[string]string
			subscription map[string]string
			account      map[string]string
		}

		fmt.Println("Auth: Enter User ID")
		fmt.Scanln(&userID)

		fmt.Println(proccessorOptions)
		fmt.Scanf("%d", &proc)

		if proc == 1 {
			processor = "card"

			data = struct {
				card         map[string]string
				subscription map[string]string
				account      map[string]string
			}{card: card}

		} else if proc == 2 {
			processor = "bank"

			data = struct {
				card         map[string]string
				subscription map[string]string
				account      map[string]string
			}{account: account}
		}

		r := Request{
			UserID: userID,
			Data:   data,
			Transaction: struct {
				transaction string
				refund      string
			}{
				transaction: payments.NewUUID(),
			},
			Processor: processor,
			Operation: r.Operation(),
		}

		if r.Operation == "ref" {
			fmt.Println("Enter refund transaction ID:")
			fmt.Scanln(&r.Transaction.refund)
		}

		return r
	}

	return Request{}
}
