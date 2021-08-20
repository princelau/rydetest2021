package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

//interface to hide all logic

type UserServiceInterface interface {
	InsertUserIntoDatabase(ctx context.Context, user *User) (string, bool, int, string)
	GetUserFromDatabase(ctx context.Context, userId string) (*User, bool, int, string)
	GetAllUsersFromDatabase(ctx context.Context) ([]User, bool, int, string)
	DeleteUserFromDatabase(ctx context.Context, userId string) (string, bool, int, string)
	UpdateUserFromDatabase(ctx context.Context, user *User, userId string) (string, bool, int, string)
}

type userService struct {
}

func NewUserService() UserServiceInterface {
	return &userService{}
}

func (c *userService) InsertUserIntoDatabase(ctx context.Context, user *User) (string, bool, int, string) {

	//if user_id is empty, return 422 error as user_id will be meaningful unique identifier
	if user.UserID == "" {
		log.Errorf("user_id cannot be empty")
		return "", false, http.StatusUnprocessableEntity, `user_id cannot be empty`
	}

	//if name is empty, return 422 error as every user requires a name
	if user.Name == "" {
		log.Errorf("name cannot be empty")
		return "", false, http.StatusUnprocessableEntity, `name cannot be empty`
	}

	collection := client.Database("database").Collection("users")

	//check if user with unique user_id already exists in the database
	var duplicateUser User
	searchCriteria := bson.M{"user_id": user.UserID}
	findErr := collection.FindOne(ctx, searchCriteria).Decode(&duplicateUser)
	//if something went wrong with finding user and it is not because no user found, return 500 error
	if findErr != nil && findErr.Error() != "mongo: no documents in result" {
		log.Errorf(findErr.Error())
		return "", false, http.StatusInternalServerError, findErr.Error()

	}

	//if successfully found an existing user in database with user_id, return 400 error
	if duplicateUser.UserID != "" {
		log.Errorf(`User with user_id: ` + duplicateUser.UserID + ` already exists in the database`)
		return "", false, http.StatusBadRequest, `User with user_id: ` + duplicateUser.UserID + ` already exists in the database`
	}

	//insert new user into database
	user.CreatedAt = time.Now().Unix()
	result, insertErr := collection.InsertOne(ctx, user)
	//if something went wrong with inserting into database, return 500 error
	if insertErr != nil {
		log.Errorf(insertErr.Error())
		return "", false, http.StatusInternalServerError, insertErr.Error()
	}

	//formatting of inserted_id
	insertedId, ok := result.InsertedID.(primitive.ObjectID)
	//if something went wrong with converting inserted id, return 500 error
	if !ok {
		log.Errorf("failed to convert inserted id into primitive object")
		return "", false, http.StatusInternalServerError, "failed to convert inserted id into primitive object"
	}

	//return insert_id if everything is successful
	log.Debugf(`Successfully inserted: ` + user.UserID + ` into database`)
	return insertedId.Hex(), true, 0, ""
}

func (c *userService) GetUserFromDatabase(ctx context.Context, userId string) (*User, bool, int, string) {

	var user *User
	collection := client.Database("database").Collection("users")
	//find user in database using user_id
	searchCriteria := bson.M{"user_id": userId}
	findErr := collection.FindOne(ctx, searchCriteria).Decode(&user)
	if findErr != nil {
		log.Errorf(findErr.Error())
		//if user does not exist in database, return 400 error
		if findErr.Error() == "mongo: no documents in result" {
			return &User{}, false, http.StatusBadRequest, findErr.Error()
		}
		//else if something went wrong with finding user, return 500 error
		return &User{}, false, http.StatusInternalServerError, findErr.Error()
	}

	//return user if everything is successful
	log.Debugf(`Successfully got user: ` + user.UserID + ` from database`)
	return user, true, 0, ""

}

func (c *userService) GetAllUsersFromDatabase(ctx context.Context) ([]User, bool, int, string) {

	var users []User
	collection := client.Database("database").Collection("users")

	//find all users from database with empty filter
	cursor, findErr := collection.Find(ctx, bson.M{})

	//if something went wrong with finding user, return 500 error
	if findErr != nil {
		log.Errorf(findErr.Error())
		return []User{}, false, http.StatusInternalServerError, findErr.Error()
	}

	//close cursor after done
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		users = append(users, user)
	}
	//if something went wrong with cursor, return 500 error
	if cursorErr := cursor.Err(); cursorErr != nil {
		log.Errorf(cursorErr.Error())
		return []User{}, false, http.StatusInternalServerError, cursorErr.Error()
	}

	//return users if everything is successful
	log.Debugf(`Successfully all users from database`)
	return users, true, 0, ""

}

func (c *userService) DeleteUserFromDatabase(ctx context.Context, userId string) (string, bool, int, string) {

	collection := client.Database("database").Collection("users")

	//find and delete user based on filter
	searchCriteria := bson.M{"user_id": userId}
	resp, deleteErr := collection.DeleteOne(ctx, searchCriteria)
	//if something went wrong during delete, return 500 error
	if deleteErr != nil {
		log.Errorf(deleteErr.Error())
		return "", false, http.StatusInternalServerError, deleteErr.Error()

	}
	//if user does not exist in database, return 400 error
	if resp.DeletedCount == 0 {
		log.Errorf(`User with user_id: ` + userId + ` does not exist in the database`)
		return "", false, http.StatusBadRequest, `User with user_id: ` + userId + ` does not exist in the database`
	}
	//return user_id of user if deletion is successful
	log.Debugf(`Successfully deleted user: ` + userId + ` from database`)
	return userId, true, 0, ""
}

func (c *userService) UpdateUserFromDatabase(ctx context.Context, user *User, userId string) (string, bool, int, string) {

	collection := client.Database("database").Collection("users")

	//if user is trying to update user_id, return 400 error as it is not allowed
	if user.UserID != userId {
		log.Errorf("user_id mismatch / tried to update with different user_id")
		return "", false, http.StatusBadRequest, "user_id mismatch / tried to update with different user_id"
	}

	// find one and update using user_id as filter and update from payload as new values
	searchCriteria := bson.M{"user_id": userId}
	_, updateErr := collection.UpdateOne(ctx, searchCriteria, bson.M{"$set": user})

	if updateErr != nil {
		log.Errorf(updateErr.Error())
		//if user does not exist in database, return 400 error
		if updateErr.Error() == "mongo: no documents in result" {
			return "", false, http.StatusBadRequest, updateErr.Error()
		} else {
			// else return 500 error if something went wrong during find or update (that is not user doesnt exist)
			return "", false, http.StatusInternalServerError, updateErr.Error()
		}
	}

	//return user_id of user if update is successful
	log.Debugf(`Successfully updated user: ` + userId + ` from database`)
	return userId, true, 0, ""

}
