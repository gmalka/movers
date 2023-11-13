package model

type Task struct {
	TaskId   int    `json:"id"`
	ItemName string `json:"name"`
	Weight   int    `json:"weight"`
}
