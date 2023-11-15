package model

type WorkerInfo struct {
	Name        string `json:"name"`
	Fatigue     int    `json:"fatigue"`
	Salary      int    `json:"salary"`
	CarryWeight int    `json:"carryweight"`
	Drunk       int    `json:"drunk"`
	Choosen     bool   `json:"choosen"`
}

type WorkerUser struct {
	Worker WorkerInfo
	User   User
}
