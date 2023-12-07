package main

import (
	"fmt"
)

func main() {
	fmt.Println("payment")

	u := NewUser("Abdul014", "John", "Trent")

	err := Save(u)
	fmt.Println(err)
}
