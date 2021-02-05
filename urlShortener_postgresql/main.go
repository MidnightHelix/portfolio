package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/MidnightHelix/urlShortener_postgresql/base62"
	"github.com/MidnightHelix/urlShortener_postgresql/models"
	"github.com/gorilla/mux"
)

type DBClient struct {
	db *sql.DB
}
type Response struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func (client *DBClient) ShortenURL(resp http.ResponseWriter, req *http.Request) {
	var id int
	var response Response
	requestBody, _ := ioutil.ReadAll(req.Body)
	error := json.Unmarshal(requestBody, &response)
	error = client.db.QueryRow("INSERT INTO urlShortener(url) VALUES($1) RETURNING id", response.URL).Scan(&id)
	responseMap := map[string]string{"encoded-string": base62.Encode(id)}

	if error != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(error.Error()))
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(responseMap)
		resp.Write(response)
	}
}

func (client *DBClient) OriginalURL(resp http.ResponseWriter, req *http.Request) {
	var url string
	vars := mux.Vars(req)
	id := base62.Decode(vars["encoded-string"])
	error := client.db.QueryRow("SELECT url FROM urlShortener WHERE id=$1", id).Scan(&url)
	if error != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(error.Error()))
	} else {
		resp.WriteHeader(http.StatusOK)
		resp.Header().Set("Content-Type", "application/json")
		responseMap := map[string]interface{}{"url": url}
		response, _ := json.Marshal(responseMap)
		resp.Write(response)
	}
}

func main() {
	db, error := models.InitDB()
	if error != nil {
		panic(error)
	}
	dbclient := &DBClient{db: db}
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/urlshortener/{encoded-string:[a-zA-Z0-9]*}", dbclient.OriginalURL).Methods("GET")
	router.HandleFunc("/urlshortener", dbclient.ShortenURL).Methods("POST")
	server := &http.Server{
		Addr:         "localhost:8000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
