package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/securecookie"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
)
var cookie = os.Setenv("SESSION_SECRET", string(securecookie.GenerateRandomKey(64)))
var secretKey = []byte(os.Getenv("SESSION_SECRET"))
var users = map[string]string{"user": "test", "admin": "password"}


type Response struct {
	Token  string `json:"token"`
	Status string `json:"status"`
}


func HealthcheckHandler(resp http.ResponseWriter, req *http.Request) {
	tokenString, err := request.HeaderExtractor{"access_token"}.ExtractToken(req)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return secretKey, nil
	})
	if err != nil {
		resp.WriteHeader(http.StatusForbidden)
		resp.Write([]byte("Access Denied; Please check the access token"))
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
	
		response := make(map[string]string)
		// response["user"] = claims["username"]
		response["time"] = time.Now().String()
		response["user"] = claims["username"].(string)
		responseJSON, _ := json.Marshal(response)
		resp.Write(responseJSON)
	} else {
		resp.WriteHeader(http.StatusForbidden)
		resp.Write([]byte(err.Error()))
	}
}


func getTokenHandler(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		http.Error(resp, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")
	if originalPassword, ok := users[username]; ok {
		if password == originalPassword {
			
			claims := jwt.MapClaims{
				"username":  username,
				"ExpiresAt": 15000,
				"IssuedAt":  time.Now().Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString(secretKey)
			if err != nil {
				resp.WriteHeader(http.StatusBadGateway)
				resp.Write([]byte(err.Error()))
			}
			response := Response{Token: tokenString, Status: "success"}
			responseJSON, _ := json.Marshal(response)
			resp.WriteHeader(http.StatusOK)
			resp.Header().Set("Content-Type", "application/json")
			resp.Write(responseJSON)

		} else {
			http.Error(resp, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(resp, "User is not found", http.StatusNotFound)
		return
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/getToken", getTokenHandler).Methods("POST")
	router.HandleFunc("/healthcheck", HealthcheckHandler).Methods("POST")
	http.Handle("/", router)
	server := &http.Server{
		Handler: router,
		Addr:    "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
