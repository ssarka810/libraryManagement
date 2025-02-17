package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/libraryManagement/config"
	"github/libraryManagement/db"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (server *Server)AddBook(w http.ResponseWriter, r *http.Request){
	bookInfo :=&config.BookDetails{}
	err :=json.NewDecoder(r.Body).Decode(bookInfo)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	err =server.ValidateStruct(bookInfo)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	dbBook,err :=server.db.GetBookByNameAndAuthor(bookInfo.Name,bookInfo.Author)
	if err!=nil || errors.Is(err,gorm.ErrRecordNotFound){
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	if dbBook!=nil{
		errStr :=fmt.Sprintf("book: [%s], author: [%s] is already present ",bookInfo.Name,bookInfo.Author)
		http.Error(w,errors.New(errStr).Error(),http.StatusBadRequest)
		return
	}
	bookInfoForDB :=server.adapter.ConvertBookDetailsForDB(bookInfo)
	bookInfoForDB.Available=true
	err =server.db.AddBooks(bookInfoForDB)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	response :=fmt.Sprintf("book: [%s], author: [%s] is added successfully ",bookInfo.Name,bookInfo.Author)
	w.Write([]byte(response))
}

func (server *Server)DeleteBook(w http.ResponseWriter, r *http.Request){
	bookInfo :=&config.BookDetails{}
	err :=json.NewDecoder(r.Body).Decode(bookInfo)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	err =server.ValidateStruct(bookInfo)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	dbBook,err :=server.db.GetBookByNameAndAuthor(bookInfo.Name,bookInfo.Author)
	if err!=nil || errors.Is(err,gorm.ErrRecordNotFound){
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	if dbBook==nil{
		errStr :=fmt.Sprintf("book: [%s], author: [%s] is not present. Nothing to delete ",bookInfo.Name,bookInfo.Author)
		http.Error(w,errors.New(errStr).Error(),http.StatusBadRequest)
		return
	}
	err =server.db.DeleteBooks(dbBook)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	response :=fmt.Sprintf("book: [%s], author: [%s] is deleted successfully ",bookInfo.Name,bookInfo.Author)
	w.Write([]byte(response))
}

func (server *Server)GetAllAvailableBooks(w http.ResponseWriter, r *http.Request){
	logrus.Info("Checking all the available books.......")
	books, err :=server.db.GetAllAvailableBook()
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	logrus.Info("books from db .. ",books)
	responseData :=&[]config.BookDetails{}
	if books!=nil{
		responseData=server.adapter.ConvertDbBooksDetailsForResponse(books)
	}
	logrus.Info("response data ",responseData)
	response,err :=json.Marshal(responseData)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (server *Server)GetBooksByName(w http.ResponseWriter, r *http.Request){
	logrus.Info("Checking all books by .......")
	name :=mux.Vars(r)["name"]
	if name==""{
		http.Error(w,errors.New("name can not be empty ").Error(),http.StatusInternalServerError)
		return
	}
	books :=&[]db.BookDetails{}
	var err error
	if  name == "all"{
		books, err =server.db.GetAllBooks()
	}else{
		books, err =server.db.GetBookByName(name)
	}
	
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	responseData :=&[]config.BookDetailsWithAvailability{}
	if books!=nil{
		responseData=server.adapter.ConvertDbBooksDetailsForResponseWithAvailability(books)
	}
	response,err :=json.Marshal(responseData)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (server *Server)GetBookByNameAndAuthor(w http.ResponseWriter, r *http.Request){
	logrus.Info("Checking all books.......")
	vars :=mux.Vars(r)
	name :=vars["name"]
	author :=vars["author"]
	if name=="" || author==""{
		http.Error(w,errors.New("name or author can not be empty ").Error(),http.StatusInternalServerError)
		return
	}
	book, err :=server.db.GetBookByNameAndAuthor(name,author)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	responseData :=&config.BookDetailsWithAvailability{}
	if book!=nil{
		responseData=server.adapter.ConvertDbBookDetailsForResponseWithAvailability(book)
	}
	response,err :=json.Marshal(responseData)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.Write(response)
}