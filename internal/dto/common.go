package dto

type PaginationRequest struct {
	Page   int    `form:"page" json:"page"`
	Limit  int    `form:"limit" json:"limit"`
	Search string `form:"search" json:"search"`
	Sort   string `form:"sort" json:"sort"`
	Dir    string `form:"dir" json:"dir"`
}
