package main

import (
	"net/http"

	"github.com/glbter/go-genesis-ses-2021/api"
	"github.com/glbter/go-genesis-ses-2021/auth"
	"github.com/glbter/go-genesis-ses-2021/dao"
	"github.com/gorilla/mux"
)

func main() {
	dao.UserDaoObj = dao.UserDao{"users.csv"}

	r := mux.NewRouter()

	r.Handle("/user/create", api.UserCreate).Methods("POST")
	r.Handle("/user/login", api.UserLogin).Methods("POST")
	r.Handle("/btcRate", auth.Authenticate(api.ExRate)).Methods("GET")

	http.ListenAndServe(":8081", r)
}
