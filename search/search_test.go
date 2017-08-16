package search

import (
	"testing"
)

var ()

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

func BenchmarkSetRadius(b *testing.B) {

}
