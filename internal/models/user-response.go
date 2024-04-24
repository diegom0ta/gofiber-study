package models

type UserResponse struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Total    int    `json:"total"`
	Data     []User `json:"data"`
}
