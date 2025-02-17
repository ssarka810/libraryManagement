package db

type UserAddress struct {
	Street string `gorm:"column:street"` 
	City string `gorm:"column:city"` 
	ZipCode string `gorm:"column:zip_code"`
}
type UserDetails struct { 
	Id int32 `gorm:"column:id;primaryKey;autoIncrement"` 
	Name string `gorm:"column:name"`
	UserName string `gorm:"colum:username;unique"`
	Role string `gorm:"column:role"`
	Password string `gorm:"column:password"` 
	Address UserAddress `gorm:"embedded"` 
	ContactNo string `gorm:"column:contactno"`
}

type User interface{
	GetUserDetails(username string)(*UserDetails,error)
	RegisterUser(userdetails *UserDetails)error
	UpdateUser(userdetails *UserDetails)error
	GetUserIdByUsername(username string)(int32,error)
}

func (store *StoreDetails)GetUserDetails(username string)(*UserDetails,error){
	user :=&UserDetails{}
	err :=store.db.Where("user_name = ?",username).First(user).Error
	if err!=nil{
		return user,err
	}
	return user,nil
}

func (store *StoreDetails)RegisterUser(userdetails *UserDetails)error{
	err :=store.db.Create(userdetails).Error
	if err!=nil{
		return err
	}
	return nil
}

func (store *StoreDetails)UpdateUser(userdetails *UserDetails)error{
	err :=store.db.Save(userdetails).Error
	if err!=nil{
		return err
	}
	return nil
}

func (store *StoreDetails)GetUserIdByUsername(username string)(int32,error){
	var userId int32
	user :=&UserDetails{}
	err :=store.db.Where("user_name = ? ",username).First(user).Error
	if err!=nil{
		return userId,err
	}
	return user.Id,nil
}