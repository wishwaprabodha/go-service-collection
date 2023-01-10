package models

import (
	"fmt"
	"math/rand"
	"strconv"
)

type Book struct {
	Id     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Name   string  `json:"bookName"`
	Author *Author `json:"author"`
}

type Author struct {
	Name  string `json:"name"`
	Genre string `json:"genre"`
}

var Books []Book

func InitBooks() {
	Books = append(Books, Book{Id: "1ddd", Isbn: "438227", Name: "Book One", Author: &Author{Name: "John", Genre: "Doe"}})
}

func ValueUpdater(newValue string, oldValue string) string {
	if newValue != "" || newValue == oldValue {
		return oldValue
	}
	return newValue
}

func (b *Book) GetBooks() []Book {
	return Books
}

func (b *Book) GetBook(id string) Book {
	for _, book := range Books {
		if id == book.Id {
			return book
		}
	}
	return Book{Id: ""}
}

func (b *Book) AddBook(newBook Book) {
	newBook.Id = strconv.Itoa(rand.Intn(100000000))
	Books = append(Books, newBook)
	// return the added book
}

func (b *Book) UpdateBook(id string, bookToUpdate Book) Book {
	for index, book := range Books {
		if book.Id == id {
			Books[index].Isbn = ValueUpdater(bookToUpdate.Isbn, Books[index].Isbn)
			Books[index].Name = ValueUpdater(bookToUpdate.Name, Books[index].Name)
			Books[index].Author.Name = ValueUpdater(bookToUpdate.Author.Name, Books[index].Author.Name)
			Books[index].Author.Genre = ValueUpdater(bookToUpdate.Author.Genre, Books[index].Author.Genre)
			Books[index] = bookToUpdate
			fmt.Println(Books)
		}
	}
	return bookToUpdate
	// de-serialize return value
}

func (b *Book) DeleteBook(id string) Book {
	var bookToDelete Book
	for index, book := range Books {
		if id == book.Id {
			bookToDelete = book
			Books = append(Books[:index], Books[index+1:]...)
		}
	}
	return bookToDelete
}
