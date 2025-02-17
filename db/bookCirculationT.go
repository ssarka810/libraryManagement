package db

import (
	"time"

	"github.com/sirupsen/logrus"
)

type BookCirculationDetails struct{
	Id int32 `gorm:"column:id;primaryKey;autoIncrement"` 
	BookId int32 `gorm:"column:book_id"` 
	UserId int32 `gorm:"column:user_id"`
	BorrowingDate time.Time `gorm:"column:borrowing_time"` 
	ReturnningDate time.Time `gorm:"column:returnning_time"` 
}

type BookCirculation interface{
	AddBookCirculationDetails(bookCirculationDetails *BookCirculationDetails)error
	UpdateBookCirculationDetails(bookCirculationDetails *BookCirculationDetails)error
	DeleteBookCirculationDetails(bookCirculationDetails *BookCirculationDetails)error
	GetBookCirculationDetailsForaUser(username string)(*[]BookCirculationDetails,error)
	GetBorrowedBookIdsByUserName(userName string) ([]int32, error) 
	GetBorrowedBooksByUserName(userName string) (*[]BookCirculationDetails, error) 
}

func (store *StoreDetails)AddBookCirculationDetails(bookCirculationDetails *BookCirculationDetails)error{
	err :=store.db.Create(bookCirculationDetails).Error
	if err!=nil{
		logrus.Info("unable to add bookCirculationDetails")
		return err
	}
	return nil
}

func (store *StoreDetails)UpdateBookCirculationDetails(bookCirculationDetails *BookCirculationDetails)error{
	err :=store.db.Save(bookCirculationDetails).Error
	if err!=nil{
		logrus.Info("unable to update bookCirculationDetails")
		return err
	}
	return nil
}

func (store *StoreDetails)DeleteBookCirculationDetails(bookCirculationDetails *BookCirculationDetails)error{
	err :=store.db.Delete(bookCirculationDetails).Error
	if err!=nil{
		logrus.Info("unable to delete bookCirculationDetails")
		return err
	}
	return nil
}

func (store *StoreDetails)GetBookCirculationDetailsForaUser(username string)(*[]BookCirculationDetails,error){
	bookCirculationDetails :=&[]BookCirculationDetails{}
	err :=store.db.Table("book_circulation_details").
	Select("book_circulation_details.*").
	Joins("JOIN user_details ON user_details.id = book_circulation_details.user_id").
	Where("user_details.user_name = ?",username).Scan(bookCirculationDetails).Error
	if err!=nil{
		return bookCirculationDetails,err
	}
	return bookCirculationDetails,nil
}

func (store *StoreDetails)GetBorrowedBooksByUserName(userName string) (*[]BookCirculationDetails, error) {
	bookCirculationDetails :=&[]BookCirculationDetails{}
	err := store.db.Table("book_circulation_details").
			Select("book_circulation_details.book_id").
			Joins("JOIN user_details ON user_details.id = book_circulation_details.user_id").
			Where("user_details.user_name = ? AND book_circulation_details.returnning_time IS NULL", userName).
			Scan(bookCirculationDetails).Error
	if err != nil {
			return nil, err
	}
	return bookCirculationDetails, nil
}

func (store *StoreDetails)GetBorrowedBookIdsByUserName(userName string) ([]int32, error) {
	var bookIDs []int32
	err := store.db.Table("book_circulation_details").
			Select("book_circulation_details.book_id").
			Joins("JOIN user_details ON user_details.id = book_circulation_details.user_id").
			Where("user_details.user_name = ? AND book_circulation_details.returnning_time IS NULL", userName).
			Pluck("book_circulation_details.book_id", &bookIDs).Error
	if err != nil {
			return nil, err
	}
	return bookIDs, nil
}

