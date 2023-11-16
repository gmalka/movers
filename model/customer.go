package model

type CustomerInfo struct {
	Name  string `json:"name"`
	Money int    `json:"money"`
	Lost  bool   `json:"lost"`
}

type CustomerUser struct {
	Customer CustomerInfo
	User     User
}
