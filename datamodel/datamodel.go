package datamodel

import (
	"fmt"
)

type Coffee struct {
	Id      string
	Name    string
	Rate    float32
	Reviews []Review
}

type Review struct {
	StoreId string
	Text    string
}

type Comment struct {
	ID      string
	PlaceID string
	Comment string
}

type Search struct {
	APIKEY   string
	KEYWORD  string
	LAT, LNG float64
}

func (s Search) Verify(data Search) bool {
	//Verify
	//check
	if len(data.KEYWORD) == 0 {
		fmt.Println("KEYWORD has wrong format !!!")
		return false
	}
	return true
}
