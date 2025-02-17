package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/libraryManagement/config"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)




func (server *Server)UserRegisterHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("Inside RegisterHandler......")
	userData :=&config.UserRegisterData{}
	err:=json.NewDecoder(r.Body).Decode(userData)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	logrus.Info("user data input : ",userData)

	err =server.ValidateStruct(userData)
	if err!=nil{
		logrus.Info("input data is invalid")
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}

	stored_db_user, err :=server.db.GetUserDetails(userData.UserName)
	if err!=nil && !errors.Is(err, gorm.ErrRecordNotFound){
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	logrus.Info("db user details ",stored_db_user)
	if stored_db_user!=nil && stored_db_user.UserName==userData.UserName{
		w.Write([]byte(`user is already present`))
		return
	}

	dbUser :=server.adapter.ConvertDataToStoreUserDetailsInDb(*userData)

	logrus.Info("user details for db operation : ",dbUser)

	err =server.db.RegisterUser(dbUser)
	if err!=nil{
		http.Error(w,err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`"user registered successfully"`))
}

func (server *Server) ValidateStruct(data interface{}) error {
	validate := validator.New()
	return validate.Struct(data)
}

func (server *Server)LoginHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("Inside LoginHandler......") 
	loginData := &config.LoginDetails{}

	err := json.NewDecoder(r.Body).Decode(loginData) 
	if err != nil { 
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Info("Input login data ......",loginData) 
	stored_db_user, err := server.db.GetUserDetails(loginData.UserName) 
	logrus.Info("stored db data ......",stored_db_user) 
	if err != nil || (stored_db_user==nil) { 
		http.Error(w, "invalid username or password", http.StatusUnauthorized) 
		return 
	}
	if (stored_db_user.Password != loginData.Password)  || (stored_db_user.UserName!=loginData.UserName){
		http.Error(w, "invalid username or password", http.StatusUnauthorized) 
		return 
	}
	session, _ := server.store.Get(r, "session-name") 
	session.Values["user"] = loginData.UserName 
	session.Options = &sessions.Options{ 
		Path: "/", 
		MaxAge: 15 * 60, // 15 minutes 
		HttpOnly: true, 
	} 
	session.Save(r, w) 
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`"login successful"`))
}

func (server *Server)UpdateAdminHandler(w http.ResponseWriter, r *http.Request){
	logrus.Info("Inside Update Admin .......")
	username :=mux.Vars(r)["username"]

	dbUser , err :=server.db.GetUserDetails(username)
	if err!=nil{
		logrus.Info("either user is not present or error while fetching data. err: ",err)
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	dbUser.Role="admin"
	err =server.db.UpdateUser(dbUser)
	if err!=nil{
		logrus.Info("error while updating user data. err: ",err)
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	response :=fmt.Sprintf("%s is an admin now",dbUser.Name)
	w.Write([]byte(response))
}