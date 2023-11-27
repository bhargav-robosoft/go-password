package db

import (
	"context"
	"net/http"
	"password-manager/constants"
	"password-manager/logger"
	"password-manager/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterUser(email string, password string) (userId string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	usersCollection := client.Database(constants.DatabaseName).Collection(constants.UsersCollection)

	insertResult, err := usersCollection.InsertOne(context.Background(), bson.M{"email": email, "password": password, "passwordSetAt": time.Now().UTC().Format(time.RFC3339)})
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func ResetPassword(email string, password string) (passwordSetAt string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	usersCollection := client.Database(constants.DatabaseName).Collection(constants.UsersCollection)

	timestamp := time.Now().UTC().Format(time.RFC3339)
	filter := bson.M{"email": email}
	update := bson.M{"$set": bson.M{"password": password, "passwordSetAt": timestamp}}
	_, err = usersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return timestamp, nil
}

func CheckPasswordReset(userId string) (passwordSetAt string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	usersCollection := client.Database(constants.DatabaseName).Collection(constants.UsersCollection)

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Invalid User Id", Status: http.StatusBadRequest}
	}

	result := usersCollection.FindOne(context.Background(), bson.M{"_id": userObjId})

	var user bson.M
	result.Decode(&user)

	return user["passwordSetAt"].(string), nil
}

func CheckUserRegistered(email string) (status bool, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return false, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	usersCollection := client.Database(constants.DatabaseName).Collection(constants.UsersCollection)

	result := usersCollection.FindOne(context.Background(), bson.M{"email": email})

	var user bson.M
	result.Decode(&user)

	if user == nil {
		return false, nil
	}

	return true, nil
}

func CheckUserCredentials(email string, password string) (status bool, userId string, passwordSetAt string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return false, "", "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	usersCollection := client.Database(constants.DatabaseName).Collection(constants.UsersCollection)

	result := usersCollection.FindOne(context.Background(), bson.M{"email": email, "password": password})

	var user bson.M
	result.Decode(&user)

	if user == nil {
		return false, "", "", nil
	}

	return true, user["_id"].(primitive.ObjectID).Hex(), user["passwordSetAt"].(string), nil
}

func CheckUserCredentialsWithId(userId string, password string) (email string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	usersCollection := client.Database(constants.DatabaseName).Collection(constants.UsersCollection)

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Invalid User Id", Status: http.StatusBadRequest}
	}

	result := usersCollection.FindOne(context.Background(), bson.M{"_id": userObjId, "password": password})

	var user bson.M
	result.Decode(&user)

	if user == nil {
		return "", nil
	}

	return user["email"].(string), nil
}

func BlacklistToken(token string, expirationTime time.Time) (err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	blacklistCollection := client.Database(constants.DatabaseName).Collection(constants.BlacklistCollection)

	// index := mongo.IndexModel{
	// 	Keys:    bson.M{"expireAt": 1},
	// 	Options: options.Index().SetExpireAfterSeconds(0),
	// }

	// _, err = blacklistCollection.Indexes().CreateOne(context.Background(), index)
	// if err != nil {
	// 	logger.ErrorLogger.Println(err.Error())
	// 	return &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	// }

	if _, err = blacklistCollection.InsertOne(context.Background(), bson.M{"token": token, "expireAt": expirationTime}); err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return nil
}

func CheckBlacklist(token string) (blacklisted bool, err error) {

	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return false, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	blacklistCollection := client.Database(constants.DatabaseName).Collection(constants.BlacklistCollection)

	result := blacklistCollection.FindOne(context.Background(), bson.M{"token": token})

	var blacklist bson.M
	result.Decode(&blacklist)

	if blacklist == nil {
		return false, nil
	}

	return true, nil
}
