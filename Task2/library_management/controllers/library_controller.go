package controllers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"library_management/models"
	"library_management/services"
)

var reader = bufio.NewReader(os.Stdin)

func Home(lib *services.Library) {
	for {
		fmt.Println("\n========== Library Management ==========")
		fmt.Println("1. Add Book")
		fmt.Println("2. Remove Book")
		fmt.Println("3. Borrow Book")
		fmt.Println("4. Return Book")
		fmt.Println("5. List Available Books")
		fmt.Println("6. List Borrowed Books")
		fmt.Println("7. Add Member")
		fmt.Println("8. Exit")

		fmt.Print("Choice: ")
		choice := readInt()

		switch choice {

		case 1:
			fmt.Print("Book ID: ")
			id := readInt()

			fmt.Print("Title: ")
			title := readString()

			fmt.Print("Author: ")
			author := readString()

			lib.AddBook(models.Book{
				ID:     id,
				Title:  title,
				Author: author,
				Status: "Available",
			})

			fmt.Println("Book added.")

		case 2:
			fmt.Print("Book ID: ")
			id := readInt()

			lib.RemoveBook(id)

			fmt.Println("Book removed.")

		case 3:
			fmt.Print("Book ID: ")
			bookID := readInt()

			fmt.Print("Member ID: ")
			memberID := readInt()

			if err := lib.BorrowBook(bookID, memberID); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book borrowed.")
			}

		case 4:
			fmt.Print("Book ID: ")
			bookID := readInt()

			fmt.Print("Member ID: ")
			memberID := readInt()

			if err := lib.ReturnBook(bookID, memberID); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Book returned.")
			}

		case 5:
			books := lib.ListAvailableBooks()

			if len(books) == 0 {
				fmt.Println("No available books.")
				continue
			}

			for _, book := range books {
				fmt.Printf("[%d] %s - %s\n",
					book.ID,
					book.Title,
					book.Author,
				)
			}

		case 6:
			fmt.Print("Member ID: ")
			memberID := readInt()

			books := lib.ListBorrowedBooks(memberID)

			if len(books) == 0 {
				fmt.Println("No borrowed books.")
				continue
			}

			for _, book := range books {
				fmt.Printf("[%d] %s - %s\n",
					book.ID,
					book.Title,
					book.Author,
				)
			}

		case 7:
			fmt.Print("Member ID: ")
			id := readInt()

			fmt.Print("Member Name: ")
			name := readString()

			err := lib.AddMember(models.Member{
				ID:   id,
				Name: name,
			})

			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Member added.")
			}

		case 8:
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}

func readString() string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readInt() int {
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		num, err := strconv.Atoi(input)
		if err == nil {
			return num
		}

		fmt.Print("Enter a valid number: ")
	}
}
