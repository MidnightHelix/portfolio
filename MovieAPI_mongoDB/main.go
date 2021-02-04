package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	collection *mongo.Collection
}

type Movie struct {
	ID        interface{} `json:"id" bson:"_id,omitempty"`
	Name      string      `json:"name" bson:"name"`
	Year      uint64      `json:"year" bson:"year"`
	Directors []string    `json:"directors" bson:"directors"`
	Writers   []string    `json:"writers" bson:"writers"`
	BoxOffice BoxOffice   `json:"boxOffice" bson:"boxOffice"`
}

type BoxOffice struct {
	Budget uint64 `json:"budget" bson:"budget"`
	Gross  uint64 `json:"gross" bson:"gross"`
}

//GET movie data
func (db *DB) GetMovie(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var movie Movie
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	error := db.collection.FindOne(context.TODO(), filter).Decode(&movie)

	if error != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(error.Error()))
	} else {
		resp.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(movie)
		resp.WriteHeader(http.StatusOK)
		resp.Write(response)
	}
}

//POST movie data
func (db *DB) PostMovie(resp http.ResponseWriter, req *http.Request) {
	var movie Movie
	requestBody, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(requestBody, &movie)
	result, error := db.collection.InsertOne(context.TODO(), movie)
	if error != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(error.Error()))
	} else {
		resp.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(result)
		resp.WriteHeader(http.StatusOK)
		resp.Write(response)
	}
}

//PUT movie data
func (db *DB) EditMovie(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	var movie Movie
	requestBody, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(requestBody, &movie)

	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": &movie}
	_, error := db.collection.UpdateOne(context.TODO(), filter, update)

	if error != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(error.Error()))
	} else {
		resp.Header().Set("Content-Type", "text/plain")
		resp.Write([]byte("Edited succesfully!"))
	}

}

//DELETE movie data
func (db *DB) DeleteMovie(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	_, error := db.collection.DeleteOne(context.TODO(), filter)
	if error != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(error.Error()))
	} else {
		resp.Header().Set("Content-Type", "text/plain")
		resp.Write([]byte("Deleted succesfully!"))
	}
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, error := mongo.Connect(context.TODO(), clientOptions)

	if error != nil {
		panic(error)
	}

	defer client.Disconnect(context.TODO())

	collection := client.Database("appDB").Collection("movies")
	db := &DB{collection: collection}

	router := mux.NewRouter()
	router.HandleFunc("/movies/{id:[a-zA-Z0-9]*}", db.GetMovie).Methods("GET")
	router.HandleFunc("/movies", db.PostMovie).Methods("POST")
	router.HandleFunc("/movies/{id:[a-zA-Z0-9]*}", db.EditMovie).Methods("PUT")
	router.HandleFunc("/movies/{id:[a-zA-Z0-9]*}", db.DeleteMovie).Methods("DELETE")

	server := &http.Server{
		Addr:         "localhost:8000",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
