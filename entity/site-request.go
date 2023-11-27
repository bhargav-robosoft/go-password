package entity

type NewSiteRequest struct {
	URL      string  `json:"url" binding:"required"`
	Name     string  `json:"name" binding:"required"`
	Sector   string  `json:"sector" binding:"required"`
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Notes    *string `json:"notes"`
}

type EditSiteRequest struct {
	Id       string  `json:"id"  binding:"required"`
	URL      string  `json:"url"`
	Name     string  `json:"name"`
	Sector   string  `json:"sector"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Notes    *string `json:"notes"`
}
