package web

import (
	"errors"
	"fmt"
	"github/libraryManagement/db"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (server *Server)BorrowBookHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("inside borrowBookHandler .........")
	vars :=mux.Vars(r)
	author :=vars["author"]
	bookName :=vars["name"]

	if author=="" ||bookName==""{
		http.Error(w, errors.New("book or author can not be empty").Error(),http.StatusBadRequest)
		return
	}
	book,err :=server.db.GetBookByNameAndAuthor(bookName,author)
	if err!=nil || book==nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	if !(book.Available){
		errStr :=fmt.Sprintf("book [%s] author[%s] is not available. Please select another book ",bookName,author)
		http.Error(w,errors.New(errStr).Error(),http.StatusInternalServerError)
		return
	}
	session, _ := server.store.Get(r, "session-name") 
	userName :=session.Values["user"].(string)
	borrowedBooks , err:=server.db.GetBorrowedBookIdsByUserName(userName)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	if len(borrowedBooks)>=3{
		errStr :="user already borrowed max number of books. return before borrowing"
		http.Error(w,errors.New(errStr).Error(),http.StatusBadRequest)
		return
	}
	userId,err :=server.db.GetUserIdByUsername(userName)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	dbBookCirculationDetails :=&db.BookCirculationDetails{
		UserId: userId,
		BookId: book.Id,
		BorrowingDate: time.Now(),
	}
	err =server.db.AddBookCirculationDetails(dbBookCirculationDetails)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	book.Available=false
	err =server.db.UpdateBooks(book)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	responseString:=fmt.Sprintf("BookName [%s], author [%s] is borrowed by [%s]",bookName,author,userName)
	w.Write([]byte(responseString))
}

func (server *Server)ReturnBookHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("inside borrowBookHandler .........")
	vars :=mux.Vars(r)
	author :=vars["author"]
	bookName :=vars["name"]

	if author=="" ||bookName==""{
		http.Error(w, errors.New("book or author can not be empty").Error(),http.StatusBadRequest)
		return
	}
	book,err :=server.db.GetBookByNameAndAuthor(bookName,author)
	if err!=nil || book==nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	session, _ := server.store.Get(r, "session-name") 
	userName :=session.Values["user"].(string)
	bookCirculationDetails , err:=server.db.GetBorrowedBooksByUserName(userName)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	ifUserBorrowed :=false
	dbBookCirculationDetails :=&db.BookCirculationDetails{}
	for _, borrowedBook :=range *bookCirculationDetails{
		if borrowedBook.BookId==book.Id{
			dbBookCirculationDetails=&borrowedBook
			ifUserBorrowed=true
			break
		}
	}
	if !ifUserBorrowed{
		http.Error(w,errors.New("user did not borrow any book. nothing to return").Error(),http.StatusInternalServerError)
		return
	}
	dbBookCirculationDetails.ReturnningDate=time.Now()
	err =server.db.UpdateBookCirculationDetails(dbBookCirculationDetails)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	book.Available=true
	err =server.db.UpdateBooks(book)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	responseString:=fmt.Sprintf("BookName [%s], author [%s] is returned by [%s]",bookName,author,userName)
	w.Write([]byte(responseString))
}