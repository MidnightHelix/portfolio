package main

import (
	"log"
	"net/http"
	"time"
	"os"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)
var cookie = os.Setenv("SESSION_SECRET", string(securecookie.GenerateRandomKey(64)))
var cookie_store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
var users = map[string]string{"user": "test", "admin": "root"}

func HealthcheckHandler(resp http.ResponseWriter, req *http.Request) {
	session, _ := cookie_store.Get(req, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		resp.Write([]byte(time.Now().String()))
	} else {
		http.Error(resp, "Forbidden", http.StatusForbidden)
	}
}


func LoginHandler(resp http.ResponseWriter, req *http.Request) {
	error := req.ParseForm()

	if error != nil {
		http.Error(resp, "Please pass the data as URL form encoded", http.StatusBadRequest)
		return
	}
	username := req.PostForm.Get("username")
	password := req.PostForm.Get("password")

	if originalPassword, ok := users[username]; ok {
		session, _ := cookie_store.Get(req, "session.id")
		if password == originalPassword {
			session.Values["authenticated"] = true
			session.Save(req, resp)
		} else {
			http.Error(resp, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(resp, "User is not found", http.StatusNotFound)
		return
	}
	resp.Write([]byte("Logged In successfully"))
	
}


func LogoutHandler(resp http.ResponseWriter, req *http.Request) {
	session, _ := cookie_store.Get(req, "session.id")
	session.Values["authenticated"] = false
	session.Save(req, resp)
	resp.Write([]byte(""))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/healthcheck", HealthcheckHandler).Methods("GET")
	router.HandleFunc("/logout", LogoutHandler).Methods("GET")
	http.Handle("/", router)

	server := &http.Server{
		Handler: router,
		Addr:    "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}