package controller

import (
	"net/http"
	"password-manager/entity"
	"password-manager/logger"
	"password-manager/service"
	"password-manager/util"

	"github.com/gin-gonic/gin"
)

type SiteController interface {
	SaveSite(ctx *gin.Context)
	GetSites(ctx *gin.Context)
	EditSite(ctx *gin.Context)
	DeleteSite(ctx *gin.Context)
}

type siteController struct {
	service service.SiteService
}

func NewSiteController(service service.SiteService) SiteController {
	return &siteController{
		service: service,
	}
}

func (controller *siteController) SaveSite(ctx *gin.Context) {
	var site entity.NewSiteRequest
	err := ctx.ShouldBindJSON(&site)
	if err != nil {
		message := "URL, Name, Sector, Username and Password are required and cannot be empty"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	userId, _ := ctx.Get("userId")

	newSite, err := controller.service.SaveSite(userId.(string), site)

	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "Site saved successfully"
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
			"site":    newSite,
		})
	}
}

func (controller *siteController) GetSites(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	sites, err := controller.service.GetSites(userId.(string))
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "Sites fetched successfully"
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
			"sites":   sites,
		})
	}
}

func (controller *siteController) EditSite(ctx *gin.Context) {
	var site entity.EditSiteRequest
	err := ctx.ShouldBindJSON(&site)
	if err != nil {
		message := "Site Id is required and cannot be empty"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	userId, _ := ctx.Get("userId")

	resultSite, err := controller.service.EditSite(userId.(string), site.Id, site)

	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "Site updated successfully"
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
			"site":    resultSite,
		})
	}
}

func (controller *siteController) DeleteSite(ctx *gin.Context) {
	siteId := ctx.Query("id")

	if siteId == "" {
		message := "Site Id is required and cannot be empty"
		logger.ErrorLogger.Println(message)
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	userId, _ := ctx.Get("userId")

	err := controller.service.DeleteSite(userId.(string), siteId)

	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "Site deleted successfully"
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
		})
	}
}
