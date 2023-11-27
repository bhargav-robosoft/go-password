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

func GenerateOtp(email string, otp string) (id string, expiresAt string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	otpCollection := client.Database(constants.DatabaseName).Collection(constants.OtpCollection)

	// index := mongo.IndexModel{
	// 	Keys:    bson.M{"expireAt": 1},
	// 	Options: options.Index().SetExpireAfterSeconds(0),
	// }

	// _, err = otpCollection.Indexes().CreateOne(context.Background(), index)
	// if err != nil {
	// 	logger.ErrorLogger.Println(err.Error())
	// 	return "", time.Time{}, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	// }

	expireTime := time.Now().UTC().Add(time.Minute * 5).Format(time.RFC3339)
	document := bson.M{
		"email":    email,
		"otp":      otp,
		"verified": false,
		"expireAt": expireTime,
	}

	result, err := otpCollection.InsertOne(context.Background(), document)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), expireTime, nil
}

func ReGenerateOtp(email string, otp string) (expiresAt string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	otpCollection := client.Database(constants.DatabaseName).Collection(constants.OtpCollection)

	expireTime := time.Now().UTC().Add(time.Minute * 5).Format(time.RFC3339)
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"otp":      otp,
			"verified": false,
			"expireAt": expireTime,
		},
	}

	_, err = otpCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return expireTime, nil
}

func VerifyOtp(dbId string, email string, otp string) (err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	otpCollection := client.Database(constants.DatabaseName).Collection(constants.OtpCollection)

	objId, err := primitive.ObjectIDFromHex(dbId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Invalid Id", Status: http.StatusBadRequest}
	}

	result := otpCollection.FindOne(context.Background(), bson.M{"_id": objId, "email": email, "otp": otp})

	var otpDocument bson.M
	result.Decode(&otpDocument)

	if otpDocument == nil {
		return &util.CustomError{Message: "Invalid OTP", Status: http.StatusBadRequest}
	} else {
		return nil
	}
}

func OtpVerified(email string) (expiresAt string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	otpCollection := client.Database(constants.DatabaseName).Collection(constants.OtpCollection)

	expireTime := time.Now().UTC().Add(time.Minute * 5).Format(time.RFC3339)
	filter := bson.M{"email": email}
	update := bson.M{
		"$set": bson.M{
			"verified": true,
			"expireAt": expireTime,
		},
	}

	if _, err = otpCollection.UpdateOne(context.Background(), filter, update); err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return expireTime, nil
}

func CheckOtpGenerated(email string) (id string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	otpCollection := client.Database(constants.DatabaseName).Collection(constants.OtpCollection)

	result := otpCollection.FindOne(context.Background(), bson.M{"email": email})

	var otpDocument bson.M
	result.Decode(&otpDocument)

	if otpDocument == nil {
		return "", nil
	} else {
		return otpDocument["_id"].(primitive.ObjectID).Hex(), nil
	}
}

func CheckOtp(email string, otp string) (err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	otpCollection := client.Database(constants.DatabaseName).Collection(constants.OtpCollection)

	expireTime := time.Now().UTC().Add(time.Minute * 5).Format(time.RFC3339)
	document := bson.M{
		"email":    email,
		"otp":      otp,
		"expireAt": expireTime,
	}

	_, err = otpCollection.InsertOne(context.Background(), document)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return nil
}

func CheckUserVerified(id string, email string) (status bool, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return false, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	otpCollection := client.Database(constants.DatabaseName).Collection(constants.OtpCollection)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return false, &util.CustomError{Message: "Invalid Id", Status: http.StatusBadRequest}
	}

	result := otpCollection.FindOne(context.Background(), bson.M{"_id": objId, "email": email, "verified": true})

	var otpDocument bson.M
	result.Decode(&otpDocument)

	if otpDocument == nil {
		return false, nil
	} else {
		return true, nil
	}
}

func RemoveVerifiedUser(id string, email string) (status bool, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return false, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	otpCollection := client.Database(constants.DatabaseName).Collection(constants.OtpCollection)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return false, &util.CustomError{Message: "Invalid Id", Status: http.StatusBadRequest}
	}

	result := otpCollection.FindOneAndDelete(context.Background(), bson.M{"_id": objId, "email": email, "verified": true})

	var otpDocument bson.M
	result.Decode(&otpDocument)

	if otpDocument == nil {
		return false, nil
	} else {
		return true, nil
	}
}
