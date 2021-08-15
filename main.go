package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

var client *mongo.Client

//main client
func main() {
	fmt.Println("Starting Application")

	//connecting to mongodb
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	//init router
	router := mux.NewRouter()

	//create controller and service
	userController := NewUserController(NewUserService())

	//apis
	router.HandleFunc("/user", userController.GetAllUsers).Methods("GET")
	router.HandleFunc("/user", userController.CreateUser).Methods("POST")
	router.HandleFunc("/user/{user_id}", userController.GetUser).Methods("GET")
	router.HandleFunc("/user/{user_id}", userController.DeleteUser).Methods("DELETE")
	router.HandleFunc("/user/{user_id}", userController.UpdateUser).Methods("PATCH")

	fmt.Println("Successfully connected to 'localhost/12345'")
	//set port address
	http.ListenAndServe(":12345", router)

}
