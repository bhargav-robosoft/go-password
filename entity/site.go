package entity

type Site struct {
	Id       string `json:"id" bson:"_id"`
	URL      string `json:"url" bson:"url"`
	Name     string `json:"name" bson:"name"`
	Sector   string `json:"sector" bson:"sector"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Notes    string `json:"notes" bson:"notes"`
	Image    string `json:"image" bson:"image"`
}

func ConvertNewSiteToSite(newSite NewSiteRequest) Site {
	return Site{
		Id:       "",
		URL:      newSite.URL,
		Name:     newSite.Name,
		Sector:   newSite.Sector,
		Username: newSite.Username,
		Password: newSite.Password,
		Notes:    *newSite.Notes,
	}
}

func ConvertEditSiteToSite(editSite EditSiteRequest, existingSite Site) Site {
	site := existingSite

	if editSite.URL != "" {
		site.URL = editSite.URL
	}
	if editSite.Name != "" {
		site.Name = editSite.Name
	}
	if editSite.Sector != "" {
		site.Sector = editSite.Sector
	}
	if editSite.Username != "" {
		site.Username = editSite.Username
	}
	if editSite.Password != "" {
		site.Password = editSite.Password
	}
	if editSite.Notes != nil {
		site.Notes = *editSite.Notes
	}

	return site
}
