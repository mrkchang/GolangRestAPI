/*
https://www.youtube.com/watch?v=SonwZ6MF5BE
go build
go build && ./restapi

http://localhost:8000/api/books
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book Stuct (Model)
type Book struct {
	ID     string  `json:"id"` // small case can be sent in from req POST
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct (Model)
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastName"`
}

// Init books var as a slice Book var
// in go, slices are variable length array. because arrays need length predefined
var books []Book

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // otherwise will be ugly string
	json.NewEncoder(w).Encode(books)
}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item) // New Encoder??
			return                          // why returning early??
		}
	}
	json.NewEncoder(w).Encode(&Book{}) // what is this '&' doing here?
}

// create books
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book) // This gets the body of the req to fills it into the book
	book.ID = strconv.Itoa(rand.Intn(1000))
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// update books
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book) // This gets the body of the req to fills it into the book
			book.ID = params["id"]
			books = append(books, book)
			// json.NewEncoder(w).Encode(book)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

// delete books
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	fmt.Println("Hello World")

	// Init router
	// typing an inference. statically typed lang
	// var age int = 35 vs. age := 35
	r := mux.NewRouter()

	// Mock data - @ todo implement DB
	books = append(books, Book{ID: "1", Isbn: "124", Title: "LOL", Author: &Author{
		Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "123123", Title: "LO2L", Author: &Author{
		Firstname: "Steve", Lastname: "Doe"}})

	// Route handler / end points
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))

}
