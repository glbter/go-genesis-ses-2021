package main

import (
	"encoding/json"
	"github.com/glbter/go-genesis-ses-2021/model"
	"github.com/glbter/go-genesis-ses-2021/util"
	"github.com/glbter/go-genesis-ses-2021/dao"
	// "fmt"
	"log"
	"io/ioutil"
	"github.com/gorilla/mux"
	"net/http"
	//"os"
	"strconv"
	"math/rand"
)

var users []model.UserLocal
var userDao dao.UserDao
func main() {
	userDao = dao.UserDao{"users.csv"}
	
	r := mux.NewRouter()

	r.Handle("/user/create", UserCreate).Methods("POST")
	r.Handle("/user/login", NotImplemented).Methods("POST")
	r.Handle("/btcRate", ExRate).Methods("GET")

	http.ListenAndServe(":8081", r)
	//log.Fatal(http.ListenAndServe(":8080", r))

	// "D:\Downloads\sem4\go-genesis-ses-2021\users.csv"
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("NotImplemented"))
})

// btcRate
var ExRate = http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {

	respData := SendRequest("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=USD&json")

	var respObjNbu []model.NbuData
	json.Unmarshal([]byte(respData), &respObjNbu)
	uahToUsd := respObjNbu[0].Rate

	respData = SendRequest("https://api.blockchain.com/v3/exchange/tickers/BTC-USD")
	
	var respObjBlockchain model.BlockchainResponse
	json.Unmarshal([]byte(respData), &respObjBlockchain)
	blockchainToUsd := respObjBlockchain.Last_trade_price

	converted := float64(blockchainToUsd) * uahToUsd
	prettyResult := util.Round(converted, 0.01)

	w.Header().Set("Content-Type", "application/json")
	payload, _ := json.Marshal(prettyResult)
	w.Write([]byte(payload))
})

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

	if userDao.GetByEmail(user.Email).Email == "" {
		log.Printf("error: the same user email: %v", user.Email)
		return
	}

	id := strconv.Itoa(rand.Intn(1000000))
	userLocal := model.UserLocal{
		Id: id,
		Name: user.Name,
		Email: user.Email,
		Password: util.Sha256(user.Password)}
	
	userDao.Create(userLocal)
})


var UserLogin = http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
	var user model.UserLogin

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	json.Unmarshal([]byte(body), &user)
})

func SendRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		//http.Error(w, "can't read body", http.StatusBadRequest)
		return nil
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	return respData
}
