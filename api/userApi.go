package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/glbter/go-genesis-ses-2021/auth"
	"github.com/glbter/go-genesis-ses-2021/dao"
	"github.com/glbter/go-genesis-ses-2021/model"
	"github.com/glbter/go-genesis-ses-2021/util"
)

// user/create
var UserCreate = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var user model.UserLogin

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	json.Unmarshal([]byte(body), &user)

	if dao.UserDaoObj.GetByEmail(user.Email).Email == "" {
		log.Printf("error: the same user email: %v", user.Email)
		http.Error(w, "user exists", http.StatusBadRequest)
		return
	}

	id := strconv.Itoa(rand.Intn(1000000))
	userLocal := model.UserLocal{
		Id:       id,
		Name:     user.Name,
		Email:    user.Email,
		Password: util.Sha256(user.Password)}

	dao.UserDaoObj.Create(userLocal)
})

// user/login
var UserLogin = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var user model.UserCredentials

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	json.Unmarshal([]byte(body), &user)

	localUsr := dao.UserDaoObj.GetByEmail(user.Email)
	if localUsr.Email == "" {
		http.Error(w, "email does not exist", http.StatusBadRequest)
		return
	}

	if util.Sha256(user.Password) != localUsr.Password {
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateJwt(localUsr.Id)

	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(token)
	w.Write([]byte(response))
})
