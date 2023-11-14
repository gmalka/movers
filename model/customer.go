package model

type CustomerInfo struct {
	Name string `json:"name"`
	Money int `json:"money"`
}

type CustomerUser struct {
	Customer CustomerInfo
	User User
}