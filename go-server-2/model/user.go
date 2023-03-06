package model

type User struct {
	UserId       int    `json:"userId"`
	UserName     string `json:"userName"`
	UserEmail    string `json:"userEmail"`
	UserPassword string `json:"userPassword"`
	ModifiedTime string `json:"modified_time"`
}
