package controllers

import (
	"fmt"
)

func Home() {
	for {
		fmt.Println("\n--- Library Management System ---")
		fmt.Println("1. Add a new book")
		fmt.Println("2. Remove an existing book")
		fmt.Println("3. Borrow a book")
		fmt.Println("4. Return a book")
		fmt.Println("5. List all available books")
		fmt.Println("6. List all borrowed books by a member")
		fmt.Println("7. Add a new member (Required to borrow books)")
		fmt.Println("8. Exit")
		fmt.Print("Choose an option: ")

	}

}
