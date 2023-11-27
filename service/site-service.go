package service

import (
	"password-manager/db"
	"password-manager/entity"
	"password-manager/logger"
	"password-manager/util"
)

type SiteService interface {
	SaveSite(userId string, site entity.NewSiteRequest) (newSite entity.Site, err error)
	GetSites(userId string) (sites []entity.Site, err error)
	EditSite(userId string, siteId string, site entity.EditSiteRequest) (resultSite entity.Site, err error)
	DeleteSite(userId string, siteId string) (err error)
}

type siteService struct{}

func NewSiteService() SiteService {
	return &siteService{}
}

func (service *siteService) SaveSite(userId string, site entity.NewSiteRequest) (newSite entity.Site, err error) {
	newSite = entity.ConvertNewSiteToSite(site)
	newSite.Image = util.GetImage(newSite.URL)
	siteId, err := db.SaveSite(userId, newSite)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return entity.Site{}, err
	}

	newSite.Id = siteId

	return newSite, nil
}

func (service *siteService) GetSites(userId string) (sites []entity.Site, err error) {
	sites, err = db.GetSites(userId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	return sites, err
}

func (service *siteService) EditSite(userId string, siteId string, updatedSite entity.EditSiteRequest) (resultSite entity.Site, err error) {
	site, err := db.GetSite(userId, siteId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return entity.Site{}, err
	}

	finalSite := entity.ConvertEditSiteToSite(updatedSite, site)
	if updatedSite.URL != site.URL {
		finalSite.Image = util.GetImage(finalSite.URL)
	}
	resultSite, err = db.EditSite(siteId, finalSite)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	return resultSite, err
}

func (service *siteService) DeleteSite(userId string, siteId string) (err error) {
	_, err = db.GetSite(userId, siteId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return err
	}

	err = db.DeleteSite(userId, siteId)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	return err
}
