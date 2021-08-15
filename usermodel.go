package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id 			primitive.ObjectID 	`json:"_id,omitempty" bson:"_id,omitempty"`
	UserID		string				`json:"user_id,omitempty" bson:"user_id,omitempty"`
	Name 		string 				`json:"name,omitempty" bson:"name,omitempty"`
	DOB 		DOB 				`json:"dob,omitempty" bson:"dob,omitempty"`
	Address 	Address 			`json:"address,omitempty" bson:"address,omitempty"`
	Description string 				`json:"description,omitempty" bson:"description,omitempty"`
	CreatedAt 	int64 				`json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type DOB struct {
	Day 	int8	`json:"day,omitempty" bson:"day,omitempty"`
	Month 	int8	`json:"month,omitempty" bson:"month,omitempty"`
	Year 	int16	`json:"year,omitempty" bson:"year,omitempty"`
}

type Address struct {
	Street 	string	`json:"street,omitempty" bson:"street,omitempty"`
	Block 	string	`json:"block,omitempty" bson:"block,omitempty"`
	Unit 	string	`json:"unit,omitempty" bson:"unit,omitempty"`
}

