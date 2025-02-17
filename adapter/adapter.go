package adapter

import (
	"github/libraryManagement/config"
	"github/libraryManagement/db"
)


type AdapterDetails struct{

}

type Adapter interface{
	ConvertDataToStoreUserDetailsInDb(config.UserRegisterData)*db.UserDetails
	ConvertBookDetailsForDB(*config.BookDetails)*db.BookDetails
	ConvertDbBookDetailsForResponse(*db.BookDetails)*config.BookDetails
	ConvertDbBooksDetailsForResponse(dbBookDetails *[]db.BookDetails)*[]config.BookDetails
	ConvertDbBookDetailsForResponseWithAvailability(dbBookDetails *db.BookDetails)*config.BookDetailsWithAvailability
	ConvertDbBooksDetailsForResponseWithAvailability(dbBookDetails *[]db.BookDetails)*[]config.BookDetailsWithAvailability
}

func NewAdapter()Adapter{
	return &AdapterDetails{}
}