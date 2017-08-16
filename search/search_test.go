package search

import (
	"testing"
)

func TestNewSearch(t *testing.T) {

}

func TestInitialize(t *testing.T) {
	searchdata := Search{}
	searchdata.apikey = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	searchdata.radius = 500
	err := Initialize(searchdata.apikey, searchdata.radius)
	if err != nil {
		t.Error(err)
	}

}

func TestPlaceSearchNoinitOK(t *testing.T) {

	var keyword string
	var lat, lng float64
	keyword = "coffee"
	lat = 25.03281
	lng = 121.33226
	apikey = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	radius = 500
	_, err := PlaceSearch(keyword, lat, lng)
	if err != ErrNotInitialized {
		t.Error(err)
	}

}

func TestPlaceSearchinitOK(t *testing.T) {

	searchdata := Search{}
	searchdata.apikey = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	searchdata.radius = 500
	err := Initialize(searchdata.apikey, searchdata.radius)
	if err != nil {
		t.Error(err)
	}

	var keyword string
	var lat, lng float64
	keyword = "coffee"
	lat = 25.03281
	lng = 121.33226
	apikey = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	radius = 500
	_, err = PlaceSearch(keyword, lat, lng)
	if err != nil {
		t.Error(err)
	}

}

func TestPlace(t *testing.T) {

	apikey = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	radius = 500
	var keyword string
	var lat, lng float64
	keyword = "coffee"
	lat = 25.03281
	lng = 121.33226

	testsearch := NewSearch(apikey, radius)
	_, err := testsearch.Place(keyword, lat, lng)
	if err != nil {
		t.Error(err)
	}

}
