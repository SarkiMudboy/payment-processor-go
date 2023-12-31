package handler

import (
	"errors"
	"fmt"
	"log"
	"payment-go/payment-processor/payments"
	"strconv"
)

type Request struct {
	UserID string
	Data   struct {
		card         map[string]string
		subscription map[string]string
		account      map[string]string
	}
	Transaction struct {
		transaction string
		refund      string
	}
	Processor string
	Operation string
}

type Handler struct {
	processor   payments.Processor
	transaction payments.Transaction
}

func (r Request) Build() Handler {

	var proc payments.Processor
	var t payments.Transaction

	switch r.Processor {
	case "card":

		data := r.Data

		// load user
		user := payments.InitUser(r.UserID)

		entries, err := payments.Load(payments.UserFile)

		if err != nil {
			log.Fatal(err)
		}

		user, err = user.Load(entries)

		// load transactions

		transaction := payments.InitTransaction(r.Transaction.transaction)

		entries, err = payments.Load(payments.TransactionFile)

		if err != nil {
			log.Fatal(err)
		}

		t, err = transaction.Load(entries)
		card := payments.NewCard(user, data.card["issuer"], data.card["number"], data.card["expiry"], data.card["cvv"])
		proc = payments.NewCreditCardProcessor(card, "Verve")
		// t = payments.NewTransaction(user, transaction.Amount, transaction.Status) i think...

	case "bank":
		data := r.Data

		// load user
		user := payments.InitUser(r.UserID)

		entries, err := payments.Load(payments.UserFile)

		if err != nil {
			log.Fatal(err)
		}

		user, err = user.Load(entries)

		// load transactions

		transaction := payments.InitTransaction(r.Transaction.transaction)

		entries, err = payments.Load(payments.TransactionFile)

		if err != nil {
			log.Fatal(err)
		}

		t, err = transaction.Load(entries)

		account := payments.NewAccount(user, data.account["holder"], data.account["number"], data.account["bank"])
		proc = payments.NewBankAccountProcessor(*account, account.Bank)
	default:
		fmt.Println("Invalid operation!")
	}

	return Handler{proc, t}
}

func (h Handler) Handle(r Request) error {

	switch r.Operation {
	case "one-time":
		h.processor.OneTimePayment(&h.transaction)
	case "sub":
		data := r.Data.subscription

		plan, err := strconv.ParseFloat(data["plan"], 64)

		if err != nil {
			log.Fatal("ERROR: ", err)
		}

		sub := payments.NewSubscription(h.transaction.User, data["name"], data["period"], plan)

		h.processor.Subscription(&sub, &h.transaction)
	case "refund":
		refundID := r.Transaction.refund
		refund := payments.InitTransaction(refundID)
		entries, err := payments.Load(payments.TransactionFile)

		if err != nil {
			log.Fatal("ERROR: ", err)
		}

		refund, _ = refund.Load(entries)

		h.processor.Refund(&h.transaction, &refund)
	default:
		return errors.New("NO OP")
	}

	return nil
}
