package datamodel

type Coffee struct {
	Id      string
	Name    string
	Rate    float32
	Reviews []Review

	TEXT []string
	Text string
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
