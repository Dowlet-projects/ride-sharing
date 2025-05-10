package models



type Make struct {
	ID int `json:"id"`
	Name string `json:"name"`
}



type Model struct {
	ID int `json:"id"`
	Name string `json:"name"`
	MakeID int `json:"make_id"`
}

