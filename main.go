package main

import (
	"encoding/json"
	"github.com/glbter/go-genesis-ses-2021/model"
	"github.com/glbter/go-genesis-ses-2021/util"
	"github.com/glbter/go-genesis-ses-2021/dao"
	"fmt"
	"log"
	"io/ioutil"
	"github.com/gorilla/mux"
	"net/http"
	//"os"
	"strconv"
	"math/rand"
	
	// "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"errors"
	"strings"
	// "github.com/gorilla/handlers"
	// "github.com/rs/cors"
	"context"
	"bytes"

	"time"
)


var users []model.UserLocal
var userDao dao.UserDao

func main() {
	userDao = dao.UserDao{"users.csv"}
	
	// jwtMidleware := jwtmiddleware.New(jwtmiddleware.Options{
	// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			
	// 		aud := "https://go-genesis-ses-2021/"
	// 		checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
	// 		if !checkAud {
	// 			return token, errors.New("Invalid audience")
	// 		}

	// 		iss := "https://dev-x9m3y4lm.eu.auth0.com/"
	// 		checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	// 		if !checkIss {
	// 			return token, errors.New("Invalid issuer")
	// 		}

	// 		cert, err := getPemCert(token)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}

	// 		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	// 		return result, nil
	// 	},

	// 	SigningMethod: jwt.SigningMethodRS256,
	// })
	// jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
	// 	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			
	// 	},
	// 	SigningMethod: jwt.SigningMethodRS256,
	// })

	
	
	r := mux.NewRouter()

	r.Handle("/user/create", UserCreate).Methods("POST")
	r.Handle("/user/login", UserLogin).Methods("POST")
	// r.Handle("/btcRate",jwtMidleware.Handler(ExRate)).Methods("GET")
	r.Handle("/btcRate", authenticate(ExRate)).Methods("GET")

	http.ListenAndServe(":8081", r)
	//log.Fatal(http.ListenAndServe(":8080", r))
}


func authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return 
		} 

		jwtToken := authHeader[1]
		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "props", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
		http.Error(w, "user exists", http.StatusBadRequest)
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

// user/login
var UserLogin = http.HandlerFunc(func(w http.ResponseWriter, r * http.Request) {
	var user model.UserLogin

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	json.Unmarshal([]byte(body), &user)

	localUsr := userDao.GetByEmail(user.Email)
	if localUsr.Email == "" {
		http.Error(w, "email does not exist", http.StatusBadRequest)
		return	
	}

	if util.Sha256(user.Password) != localUsr.Password {
		http.Error(w, "wrong password", http.StatusBadRequest)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// response, _ := json.Marshal(getJwtToken(localUsr.Id))
	
	// fmt.Println(response)
	// w.Write([]byte(response))
	// var creds Credentials
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Id: localUsr.Id,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(tokenString)
	
	w.Write([]byte(response))
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

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://dev-x9m3y4lm.eu.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func getJwtToken(userId string) AuthOResponse {

	secret := "153qPHF8ZyZaPfMlZlQvoNGdbGRd4tOP26g7LSBKLKjqFvd1zIUPrTChbk8QZAQ"
	url := "https://dev-x9m3y4lm.eu.auth0.com/oauth/token"
	aud := "https://go-genesis-ses-2021/"
	
	payload := fmt.Sprintf("{\"client_id\":\"%v\",\"client_secret\":\"%v\",\"audience\":\"%v\",\"grant_type\":\"client_credentials\"}",
		userId,
		secret,
		aud)

	fmt.Println(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(body)
	var token AuthOResponse
	json.Unmarshal([]byte(body), &token)
	fmt.Println(token)
	return token
}

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JsonWebKeys `json:"keys"`
}

type JsonWebKeys struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N string `json:"n"`
	E string `json:"e"`
	X5c []string `json:"x5c"`
}

type AuthOResponse struct {
	AccessToken string `json:"access_token"`
	tokenType string `json:"token_type"`
}







// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")


// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type AuthHeader struct {
	Authentication string `json:"authentication"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Id string `json:"id"`
	jwt.StandardClaims
}