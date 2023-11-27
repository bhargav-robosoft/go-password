package db

import (
	"context"
	"net/http"
	"password-manager/constants"
	"password-manager/entity"
	"password-manager/logger"
	"password-manager/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ReadSites() (sites []entity.Site, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return []entity.Site{}, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	sitesCollection := client.Database(constants.DatabaseName).Collection(constants.SitesCollection)

	cursor, err := sitesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return []entity.Site{}, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	sites = []entity.Site{}
	for cursor.Next(context.Background()) {
		var site entity.Site
		cursor.Decode(&site)
		sites = append(sites, site)
	}

	return sites, nil
}

func SaveSite(userId string, site entity.Site) (id string, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	sitesCollection := client.Database(constants.DatabaseName).Collection(constants.SitesCollection)

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Invalid User Id", Status: http.StatusBadRequest}
	}

	document := bson.M{
		"userId":   userObjId,
		"url":      site.URL,
		"name":     site.Name,
		"sector":   site.Sector,
		"username": site.Username,
		"password": site.Password,
		"notes":    site.Notes,
		"image":    site.Image,
	}
	result, err := sitesCollection.InsertOne(context.Background(), document)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func GetSites(userId string) (sites []entity.Site, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return []entity.Site{}, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	sitesCollection := client.Database(constants.DatabaseName).Collection(constants.SitesCollection)

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return []entity.Site{}, &util.CustomError{Message: "Invalid User Id", Status: http.StatusBadRequest}
	}

	cursor, err := sitesCollection.Find(context.Background(), bson.M{"userId": userObjId})
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return []entity.Site{}, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	sites = []entity.Site{}
	for cursor.Next(context.Background()) {
		var site entity.Site
		err = cursor.Decode(&site)
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			return []entity.Site{}, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
		}
		sites = append(sites, site)
	}

	return sites, nil
}

func EditSite(siteId string, site entity.Site) (updatedSite entity.Site, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return entity.Site{}, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	sitesCollection := client.Database(constants.DatabaseName).Collection(constants.SitesCollection)

	siteObjId, err := primitive.ObjectIDFromHex(siteId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return entity.Site{}, &util.CustomError{Message: "Invalid Site Id", Status: http.StatusBadRequest}
	}

	filter := bson.M{"_id": siteObjId}
	update := bson.M{"$set": bson.M{
		"url":      site.URL,
		"name":     site.Name,
		"sector":   site.Sector,
		"username": site.Username,
		"password": site.Password,
		"notes":    site.Notes,
		"image":    site.Image,
	}}
	options := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := sitesCollection.FindOneAndUpdate(context.Background(), filter, update, options)

	result.Decode(&updatedSite)
	return updatedSite, nil
}

func DeleteSite(userId string, siteId string) (err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	sitesCollection := client.Database(constants.DatabaseName).Collection(constants.SitesCollection)

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Invalid User Id", Status: http.StatusBadRequest}
	}
	siteObjId, err := primitive.ObjectIDFromHex(siteId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Invalid Site Id", Status: http.StatusBadRequest}
	}

	filter := bson.M{
		"_id":    siteObjId,
		"userId": userObjId,
	}
	update := bson.M{"$set": bson.M{
		"userId":    "",
		"oldUserId": userObjId,
	}}
	_, err = sitesCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	} else {
		return nil
	}
}

func GetSite(userId string, siteId string) (site entity.Site, err error) {
	client, err := DbSetup()
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return entity.Site{}, &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}
	defer client.Disconnect(context.Background())

	sitesCollection := client.Database(constants.DatabaseName).Collection(constants.SitesCollection)

	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return entity.Site{}, &util.CustomError{Message: "Invalid User Id", Status: http.StatusBadRequest}
	}
	siteObjId, err := primitive.ObjectIDFromHex(siteId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return entity.Site{}, &util.CustomError{Message: "Invalid Site Id", Status: http.StatusBadRequest}
	}

	result := sitesCollection.FindOne(context.Background(), bson.M{"_id": siteObjId, "userId": userObjId})

	var siteDecoded bson.M
	result.Decode(&siteDecoded)
	if siteDecoded == nil {
		return entity.Site{}, &util.CustomError{Message: "Site not found", Status: http.StatusBadRequest}
	}

	result.Decode(&site)

	return site, nil
}
