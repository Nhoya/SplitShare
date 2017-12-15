package main

import (
	"fmt"
	"os"

	"github.com/SSSaaS/sssa-golang"
)

func main() {
	var n int
	var shares []string
	fmt.Print("Number of slices: ")
	fmt.Scan(&n)

	if n < 1 {
		fmt.Println("You can't inster a negative integer")
		os.Exit(1)
	}
	for i := 0; i < n; i++ {
		var share string
		fmt.Print("Insert share :")
		fmt.Scanln(&share)
		shares = append(shares, share)
	}
	secret, err := sssa.Combine(shares)
	if err != nil {
		fmt.Println("Error during Decryption")
		os.Exit(1)
	}
	fmt.Println(secret)

}
