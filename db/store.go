package db

import "github.com/jinzhu/gorm"

type StoreDetails struct{
	db *gorm.DB
	UserDetails
	BookDetails
	BookCirculationDetails
}

type Store interface{
	User
	Book
	BookCirculation
}

func NewDb(db *gorm.DB)Store{
	return &StoreDetails{
		db: db,
	}
}