package main

import (
	"fmt"
	"github/libraryManagement/adapter"
	"github/libraryManagement/db"
	"github/libraryManagement/web"
)

func main(){
	fmt.Println("starting the server...")
	database :=db.Init()
	store :=db.NewDb(database)
	adapter :=adapter.NewAdapter()
	server :=web.NewServer(store,adapter)
	server.Init()
	server.Start("9090")
}