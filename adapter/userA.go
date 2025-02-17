package adapter

import (
	"github/libraryManagement/config"
	"github/libraryManagement/db"
	"strconv"
)


func (adapter *AdapterDetails)ConvertDataToStoreUserDetailsInDb(userData config.UserRegisterData)*db.UserDetails{
	dbUser :=&db.UserDetails{}
	dbUser.Name=userData.Name
	dbUser.UserName=userData.UserName
	dbUser.Password=userData.Password
	dbUser.ContactNo=userData.ContactNo
	dbUser.Address.City=userData.Address.City
	dbUser.Address.Street=userData.Address.Street
	dbUser.Address.ZipCode=strconv.Itoa(int(userData.Address.PinCode))
	dbUser.Role="user"
	return dbUser

}