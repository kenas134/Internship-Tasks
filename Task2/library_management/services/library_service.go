package services

import (
	"errors"
	"library_management/models" // Replace 'library_management' with the module name from your go.mod
)

type LibraryManager interface {
	AddBook(book models.Book)
	RemoveBook(bookID int)
	BorrowBook(bookID int, memberID int) error
	ReturnBook(bookID int, memberID int) error
	ListAvailableBooks() []models.Book
	ListBorrowedBooks(memberID int) []models.Book
}

type Library struct {
	Books   map[int]models.Book
	Members map[int]models.Member
}

// constructor
func NewLibrary() *Library {
	return &Library{
		Books:   make(map[int]models.Book),
		Members: make(map[int]models.Member),
	}
}

//Add Book

func (l *Library) AddBook(book models.Book) {
	l.Books[book.ID] = book
}

//delete Book

func (l *Library) RemoveBook(bookID int) {
	delete(l.Books, bookID)
}

//Borrow Book
// errors that can occur
//1.no book
//2.no member
//3.book already borrowed
//oprations we perform
//1.change status of the book
//2.add to the persons borrowed list the book

func (l *Library) BorrowBook(bookID int, memberID int) error {
	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("book does not exist")
	}

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("member does not exist")
	}

	if book.Status == "Borrowed" {
		return errors.New("book is already borrowed")
	}

	book.Status = "Borrowed"
	l.Books[bookID] = book

	member.BorrowedBooks = append(member.BorrowedBooks, book)
	l.Members[memberID] = member

	return nil
}

//ReturnBook
//errors to take care
//1.no book
//2.no member
//3.book already available
//operations
//1.change status to available
//2.remove it from slice

func (l *Library) ReturnBook(bookID int, memberID int) error {
	book, ok := l.Books[bookID]
	if !ok {
		return errors.New("book does not exist")
	}

	member, ok := l.Members[memberID]
	if !ok {
		return errors.New("member does not exist")
	}

	if book.Status == "Available" {
		return errors.New("book is not Borrowed")
	}

	found := false
	arr := []models.Book{}
	for _, val := range member.BorrowedBooks {
		if val.ID != bookID {
			arr = append(arr, val)
		} else {
			found = true
		}
	}
	if !found {
		return errors.New("member did not borrow the book")
	}
	book.Status = "Available"
	l.Books[bookID] = book

	member.BorrowedBooks = arr
	l.Members[memberID] = member
	return nil
}

//ListAvailableBooks
//errors
//operation
//1.filter book by checking their status

func (l *Library) ListAvailableBooks() []models.Book {

	arr := []models.Book{}

	for _, val := range l.Books {
		if val.Status == "Available" {
			arr = append(arr, val)
		}
	}
	return arr
}

//ListBorrowedBooks
//errors
//oprations
//1.filter book by checking their status

func (l *Library) ListBorrowedBooks(memberID int) []models.Book {
	return l.Members[memberID].BorrowedBooks
}

//ADD member
//errors member already exist
//operation
//1.insert member object into the member map

func (l *Library) AddMember(member models.Member) error {
	if _, exists := l.Members[member.ID]; exists {
		return errors.New("member already exists")
	}

	l.Members[member.ID] = member
	return nil
}
