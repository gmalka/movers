package model

type MoverInfo struct {
	Name        string `json:"name"`
	Fatigue     int    `json:"fatigue"`
	Salary      int    `json:"salary"`
	CarryWeight int    `json:"carryweight"`
	Drunk       int    `json:"drunk"`
}
