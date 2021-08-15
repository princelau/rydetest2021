package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

//controller layer where all the APIs are defined

type UserController struct {
	UserService UserServiceInterface
}

func NewUserController(userService UserServiceInterface) *UserController {
	return &UserController{
		userService,
	}
}

func (c *UserController) CreateUser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//decode user payload
	var user User
	decodeErr := json.NewDecoder(request.Body).Decode(&user)
	if decodeErr != nil {
		DisplayError(decodeErr.Error(), http.StatusBadRequest, response)
		return
	}

	//calling service layer to insert user into database
	result, ok, errorCode, errorMessage := c.UserService.InsertUserIntoDatabase(ctx, &user)
	if !ok {
		DisplayError(errorMessage, errorCode, response)
		return
	}

	//give a meaningful response of insert id
	response.Write([]byte(`{ "message": "InsertID :` + result + `"}`))
}

func (c *UserController) GetUser(response http.ResponseWriter, request *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	response.Header().Add("content-type", "application/json")

	//getting path params
	params := mux.Vars(request)
	userId, _ := params["user_id"]

	//calling service layer to insert user into database
	result, ok, errorCode, errorMessage := c.UserService.GetUserFromDatabase(ctx, userId)
	if !ok {
		DisplayError(errorMessage, errorCode, response)
		return
	}

	//return user struct
	json.NewEncoder(response).Encode(result)
}

func (c *UserController) DeleteUser(response http.ResponseWriter, request *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	response.Header().Add("content-type", "application/json")

	//getting path params
	params := mux.Vars(request)
	userId, _ := params["user_id"]

	//calling service layer to delete user from database
	result, ok, errorCode, errorMessage := c.UserService.DeleteUserFromDatabase(ctx, userId)
	if !ok {
		DisplayError(errorMessage, errorCode, response)
		return
	}
	//give a meaningful response of user id that has been deleted
	response.Write([]byte(`{ "message": "User with user_id: ` + result + ` has been successfully deleted from the database"}`))
}

func (c *UserController) UpdateUser(response http.ResponseWriter, request *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	response.Header().Add("content-type", "application/json")

	//decode user payload
	var user User
	decodeErr := json.NewDecoder(request.Body).Decode(&user)
	if decodeErr != nil {
		DisplayError(decodeErr.Error(), http.StatusBadRequest, response)
		return
	}

	//getting path params
	params := mux.Vars(request)
	userId, _ := params["user_id"]

	//calling service layer to delete user from database
	result, ok, errorCode, errorMessage := c.UserService.UpdateUserFromDatabase(ctx, &user, userId)
	if !ok {
		DisplayError(errorMessage, errorCode, response)
		return
	}
	//give a meaningful response of user id that has been updated
	response.Write([]byte(`{ "message": "User with user_id: ` + result + ` has been successfully updated"}`))
}


func (c *UserController) GetAllUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//calling service layer to get all users from database
	result, ok, errorCode, errorMessage := c.UserService.GetAllUsersFromDatabase(ctx)
	if !ok {
		DisplayError(errorMessage, errorCode, response)
		return
	}

	//all users from database
	json.NewEncoder(response).Encode(result)
}
