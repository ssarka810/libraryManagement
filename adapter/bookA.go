package adapter

import (
	"github/libraryManagement/config"
	"github/libraryManagement/db"
)

func (adapter *AdapterDetails)ConvertBookDetailsForDB(inputBookDetails *config.BookDetails)*db.BookDetails{
	dbBookDetails :=&db.BookDetails{}
	dbBookDetails.Author=inputBookDetails.Author
	dbBookDetails.Name=inputBookDetails.Name

	return dbBookDetails
}

func (adapter *AdapterDetails)ConvertDbBookDetailsForResponse(dbBookDetails *db.BookDetails)*config.BookDetails{
	responseBookDetails :=&config.BookDetails{}
	responseBookDetails.Author=dbBookDetails.Author
	responseBookDetails.Name=dbBookDetails.Name

	return responseBookDetails
}
func (adapter *AdapterDetails)ConvertDbBookDetailsForResponseWithAvailability(dbBookDetails *db.BookDetails)*config.BookDetailsWithAvailability{
	responseBookDetails :=&config.BookDetailsWithAvailability{}
	responseBookDetails.Author=dbBookDetails.Author
	responseBookDetails.Name=dbBookDetails.Name

	return responseBookDetails
}

func (adapter *AdapterDetails)ConvertDbBooksDetailsForResponse(dbBookDetails *[]db.BookDetails)*[]config.BookDetails{
	responseBooksDetails :=[]config.BookDetails{}
	for _,resp :=range *dbBookDetails{
		bookDetails :=config.BookDetails{}
		bookDetails.Author=resp.Author
		bookDetails.Name=resp.Name
		responseBooksDetails=append(responseBooksDetails,bookDetails)
	}
	return &responseBooksDetails
}

func (adapter *AdapterDetails)ConvertDbBooksDetailsForResponseWithAvailability(dbBookDetails *[]db.BookDetails)*[]config.BookDetailsWithAvailability{
	responseBooksDetails :=[]config.BookDetailsWithAvailability{}
	for _,resp :=range *dbBookDetails{
		bookDetails :=config.BookDetailsWithAvailability{}
		bookDetails.Author=resp.Author
		bookDetails.Name=resp.Name
		bookDetails.Available=resp.Available
		responseBooksDetails=append(responseBooksDetails,bookDetails)
	}
	return &responseBooksDetails
}