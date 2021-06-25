package main

import (
	"net/http"

	"github.com/glbter/go-genesis-ses-2021/api"
	"github.com/glbter/go-genesis-ses-2021/auth"
	"github.com/glbter/go-genesis-ses-2021/config"
	"github.com/glbter/go-genesis-ses-2021/dao"
	"github.com/gorilla/mux"
)

func main() {

	conf, err := config.InitConfig()
	if err != nil {
		return
	}

	dao.UserDaoObj = dao.UserDao{conf.DbName}

	r := mux.NewRouter()

	r.Handle("/user/create", api.UserCreate).Methods("POST")
	r.Handle("/user/login", api.UserLogin).Methods("POST")
	r.Handle("/btcRate", auth.Authenticate(api.ExRate)).Methods("GET")

	http.ListenAndServe(":8081", r)
}
