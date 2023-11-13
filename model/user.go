package model

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserInfo struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type Tokens struct {
	AccessToken string `json:"access"`
	RefreshToken string `json:"refresh"`
}