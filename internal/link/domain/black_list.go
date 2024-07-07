package domain

type BlackList struct {
	IP     string `json:"ip"`
	Status int    `json:"status"`
}

type BlackLists []BlackList
