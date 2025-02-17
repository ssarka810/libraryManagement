package db

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type BookDetails struct { 
	Id int32 `gorm:"column:id;primaryKey;autoIncrement"` 
	Name string `gorm:"column:name"` 
	Author string `gorm:"column:author"` 
	Available bool `gorm:"column:available"`
}

type Book interface{
	GetBookByName(bookName string)(*[]BookDetails,error)
	GetBookByNameAndAuthor(bookName, author string)(*BookDetails,error)
	GetAllAvailableBook()(*[]BookDetails,error)
	GetAllBooks()(*[]BookDetails,error)
	AddBooks(bookDetails *BookDetails)error
	DeleteBooks(bookDetails *BookDetails)error
	UpdateBooks(bookDetails *BookDetails)error
	GetBooksByAuthor(authorname string)(*[]BookDetails,error)
}

func (store *StoreDetails)GetBookByName(bookName string)(*[]BookDetails,error){
	books :=&[]BookDetails{}
	err :=store.db.Where("name =?",bookName).Find(books).Error
	if err!=nil{
		if errors.Is(err, gorm.ErrRecordNotFound) { 
			return nil, nil 
		}
		return books,err
	}
	return books,nil
}
func (store *StoreDetails)GetBooksByAuthor(authorname string)(*[]BookDetails,error){
	books :=&[]BookDetails{}
	err :=store.db.Where("author =?",authorname).Find(books).Error
	if err!=nil{
		if errors.Is(err, gorm.ErrRecordNotFound) { 
			return nil, nil 
		}
		return books,err
	}
	return books,nil
}
func (store *StoreDetails)GetBookByNameAndAuthor(bookName, author string)(*BookDetails,error){
	book :=&BookDetails{}
	err :=store.db.Where("name =? AND author =?",bookName,author).First(book).Error
	if err!=nil{
		if errors.Is(err, gorm.ErrRecordNotFound) { 
			return nil, nil 
		}
		return book,err
	}
	return book,nil
}

func (store *StoreDetails)GetAllAvailableBook()(*[]BookDetails,error){
	books :=&[]BookDetails{}
	err :=store.db.Where("available =?",true).Find(books).Error
	if err!=nil{
		if errors.Is(err, gorm.ErrRecordNotFound) { 
			return nil, nil 
		}
		return books,err
	}
	return books,nil
}

func (store *StoreDetails)GetAllBooks()(*[]BookDetails,error){
	books :=&[]BookDetails{}
	err :=store.db.Find(books).Error
	if err!=nil{
		return books,err
	}
	return books,nil
}

func (store *StoreDetails)AddBooks(bookDetails *BookDetails)error{
	err :=store.db.Create(bookDetails).Error
	if err!=nil{
		return err
	}
	return nil
}

func (store *StoreDetails)DeleteBooks(bookDetails *BookDetails)error{
	err :=store.db.Delete(bookDetails).Error
	if err!=nil{
		return err
	}
	return nil
}

func (store *StoreDetails)UpdateBooks(bookDetails *BookDetails)error{
	err :=store.db.Save(bookDetails).Error
	if err!=nil{
		return err
	}
	return nil
}

