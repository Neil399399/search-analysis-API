package search

import (
	"search-analysis-API/datamodel"
	"testing"
)

/*
func TestNewSearch(t *testing.T) {
	var testradius uint
	testkey := "afjadkf"
	testradius = 500
	newsearchresult := NewSearch(testkey, testradius)
	if newsearchresult.apikey != testkey {
		t.Error(ErrNoAPIKey)
	}
	if newsearchresult.radius != testradius {
		t.Error(ErrNoRadius)
	}

}
func TestInitialize(t *testing.T) {

	searchdata := Search{}
	searchdata.apikey = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	searchdata.radius = 500
	err := Initialize(searchdata.apikey, searchdata.radius)
	if err != nil {
		t.Error(err)
	}
	if initOk != true {
		t.Error(err)
	}

}
func TestSetRadius(t *testing.T) {
	var key string
	var rad uint
	S := NewSearch(key, rad)
	err := S.SetRadius(500)
	if err != nil {
		t.Error(err)
	}

}

//Benchmark
*/
func BenchmarkNewSearch(b *testing.B) {
	// run the Fib function b.N times
	var testradius uint
	testkey := "afjadkf"
	testradius = 500

	for n := 0; n < b.N; n++ {
		_ = NewSearch(testkey, testradius)

	}

}

func BenchmarkInitialize(b *testing.B) {
	// run the Fib function b.N times
	searchdata := Search{}
	searchdata.apikey = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	searchdata.radius = 500

	for n := 0; n < b.N; n++ {
		_ = Initialize(searchdata.apikey, searchdata.radius)
	}

}

func BenchmarkSetRadius(b *testing.B) {
	// run the Fib function b.N times
	var key string
	var rad uint
	S := NewSearch(key, rad)
	for n := 0; n < b.N; n++ {
		_ = S.SetRadius(500)
	}

}

func BenchmarkPlaceSearch(b *testing.B) {
	// run the Fib function b.N times
	var key, keyword string
	var rad uint
	var lat, lng float64
	key = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	rad = 500
	keyword = "海鮮餐廳"
	lat = 25.03978
	lng = 121.548495
	S := NewSearch(key, rad)
	for n := 0; n < b.N; n++ {
		_, _ = S.Place(keyword, lat, lng)
	}
}

func BenchmarkPlaceSearchLayer1(b *testing.B) {
	// run the Fib function b.N times
	var key, keyword string
	var rad uint
	var lat, lng float64
	key = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	rad = 500
	keyword = "海鮮餐廳"
	lat = 25.03978
	lng = 121.548495
	S := NewSearch(key, rad)
	for n := 0; n < b.N; n++ {
		_ = S.Placesearchlayer1(keyword, lat, lng)
	}

}

func BenchmarkPlaceSearchLayer2(b *testing.B) {
	// run the Fib function b.N times
	var key, keyword string
	var rad uint
	var lat, lng float64
	key = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	rad = 500
	keyword = "海鮮餐廳"
	lat = 25.03978
	lng = 121.548495
	S := NewSearch(key, rad)
	for n := 0; n < b.N; n++ {
		_ = S.Placesearchlayer1(keyword, lat, lng)
	}

}

func BenchmarkSearchResultWriteIndatamodel(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		coffeeList := []datamodel.Coffee{}

		for i := 0; i < 20; i++ {

			cof := datamodel.Coffee{}
			cof.Id = "asdfadf"
			cof.Name = "hello"
			cof.Rate = 9.9
			cof.Reviews = []datamodel.Review{}

			for j := 0; j < 5; j++ {

				review := datamodel.Review{cof.Id, "fdhghdghasdfadf"}
				cof.Reviews = append(cof.Reviews, review)
			}

			coffeeList = append(coffeeList, cof)
		}

	}

}
