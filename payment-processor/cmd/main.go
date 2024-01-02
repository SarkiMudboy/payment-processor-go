package main

import (
	"fmt"
	"log"
	"os"
)

func intro() {

	fmt.Println(`
	 ________________
	|				 |
	|	GO-WALLET    |
	|________________|

<WELCOME TO GO-WALLET>

	`)
}

func main() {

	intro()

	if err := root(os.Args[1:]); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
