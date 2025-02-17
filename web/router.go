package web

import (
	"fmt"
	"github/libraryManagement/adapter"
	"github/libraryManagement/db"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type Server struct{
	db db.Store
	router *mux.Router
	adapter adapter.Adapter
	store *sessions.CookieStore
}

func NewServer(db db.Store, adapter adapter.Adapter)Server{
	return Server{
		db: db,
		adapter: adapter,
	}
}


func (server *Server)Start(port string){
	fmt.Println("starting server with port ",port)
	if err :=http.ListenAndServe(fmt.Sprintf(":%s",port),server.router);err!=nil{
		logrus.Error("error while starting server ",err)
	}
}

func (server *Server) Init() {
	server.store = sessions.NewCookieStore([]byte("secret-key"))
	server.setupRouter()
}


func (server *Server)setupRouter(){
	server.router=mux.NewRouter().StrictSlash(true)
	server.addApis()
}

func (server *Server)addApis(){
	router :=server.router
	router.Use(server.SessionMiddleware)
	router.HandleFunc("/",func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`hello world`))
	})
	router.HandleFunc("/v1/user/register",server.UserRegisterHandler).Methods("POST")
	router.HandleFunc("/v1/user/login", server.LoginHandler).Methods("POST") 
	router.PathPrefix("/v1/user/{username}/admin").Handler(server.AdminMiddleware(http.HandlerFunc(server.UpdateAdminHandler)))

	//books
	router.PathPrefix("/v1/book/add").Handler(server.AdminMiddleware(http.HandlerFunc(server.AddBook))).Methods("POST")
	router.PathPrefix("/v1/book/delete").Handler(server.AdminMiddleware(http.HandlerFunc(server.DeleteBook))).Methods("DELETE")
	router.HandleFunc("/v1/book/all/available",server.GetAllAvailableBooks).Methods("GET")
	router.HandleFunc("/v1/book/{name}",server.GetBooksByName).Methods("GET")
	router.HandleFunc("/v1/book/{name}/{author}",server.GetBookByNameAndAuthor).Methods("GET")

	//bookCirculation
	router.HandleFunc("/v1/book/borrow/{name}/{author}",server.BorrowBookHandler).Methods("POST")
	router.HandleFunc("/v1/book/return/{name}/{author}",server.ReturnBookHandler).Methods("PUT")
}

// SessionMiddleware validates the session 
func (server *Server) SessionMiddleware(next http.Handler) http.Handler { 
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { 
		if r.URL.Path == "/v1/user/register" || r.URL.Path == "/v1/user/login" { 
			next.ServeHTTP(w, r) 
			return 
		} 
		session, _ := server.store.Get(r, "session-name") 
		if session.Values["user"] == nil { 
			http.Error(w, "unauthorized", http.StatusUnauthorized) 
			return 
		} 
		next.ServeHTTP(w, r) 
	})
}

func (server *Server) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := server.store.Get(r, "session-name")
			username, ok := session.Values["user"].(string)
			if !ok {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
			}

			user, err := server.db.GetUserDetails(username)
			if err != nil || user == nil {
					http.Error(w, "unauthorized", http.StatusUnauthorized)
					return
			}

			if user.Role != "admin" {
					http.Error(w, "forbidden", http.StatusForbidden)
					return
			}

			next.ServeHTTP(w, r)
	})
}
