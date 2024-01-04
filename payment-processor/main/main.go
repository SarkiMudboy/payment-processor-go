package main

import (
	"fmt"
	"log"
	"payment-go/payment-processor/payments"
)

func main() {

	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("Error: File path error")
	// }
	// fmt.Println(dir)

	fmt.Println("Hello welcome to payments!")

	user := payments.NewUser("Yahya322", "Usman", "Yahya")
	err := payments.Save(&payments.UserFile, user)

	if err != nil {
		log.Fatal(err)
	}

	// u := payments.InitUser("2e14072a-ba0e-4594-b393-ead3ea7a48da")

	// file, err := payments.Load(payments.UserFile)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// u, err = u.Load(file)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// card := payments.NewCard(u, "Verve", "345-5689-90988-33", "2030-05-04 03:04:05", "454")
	// card.Balance = 1500.32
	// card.Limit = 100000000.00
	// t := payments.NewTransaction(u, 344.90, "pending")

	// label := "Chase"

	// processor := payments.NewCreditCardProcessor(card, label)
	// processor.OneTimePayment(t)

	// err = payments.Save(&payments.TransactionFile, t)
	// if err != nil {
	// 	fmt.Println(err)
	// }

}
