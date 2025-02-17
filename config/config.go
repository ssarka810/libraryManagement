package config

type UserAddress struct {
	Street  string `json:"street,omitempty" validate:"required,min=3"`
	City    string `json:"city,omitempty" validate:"required,min=3"`
	PinCode int32  `json:"pincode,omitempty" validate:"required"`
}

type UserRegisterData struct {
	Name      string      `json:"name,omitempty" validate:"required,min=2"`
	UserName  string      `json:"username,omitempty" validate:"required,min=6"`
	Password  string      `json:"password,omitempty" validate:"required,min=6"`
	Address   UserAddress `json:"useraddress,omitempty" validate:"required"`
	ContactNo string      `json:"contactNo,omitempty" validate:"required,len=10"`
}

type LoginDetails struct{
	UserName  string      `json:"username,omitempty" validate:"required,min=6"`
	Password  string      `json:"password,omitempty" validate:"required,min=6"`
}

type BookDetails struct{
	Name      string      `json:"name,omitempty" validate:"required,min=2"`
	Author  string      `json:"author,omitempty" validate:"required,min=2"`
}

type BookDetailsWithAvailability struct{
	Name      string      `json:"name,omitempty" validate:"required,min=2"`
	Author  string      `json:"author,omitempty" validate:"required,min=2"`
	Available  bool      `json:"available,omitempty"`
}