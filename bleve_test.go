package main

import (
	"search-analysis-API/datamodel"
	"testing"
)

func Testinit(t *testing.T) {

}

func Testjieba(t *testing.T) {

	testmodel := []datamodel.Coffee{}
	testmodel[0].Id = "ab23bc888a##$%^bc"
	testmodel[0].Name = "hellocoffee"
	testmodel[0].Rate = 9.9
	testmodel[0].Reviews[0].StoreId = "ab23bc888a##$%^bc"
	testmodel[0].Reviews[0].Text = "內裝舒適，座位寛敞，聊天小歇的好地方。"

	querys = []string{
		"舒適",
		"好地方",
		"好",
	}
	_, err := jiebatest(testmodel, querys)
	if err != nil {
		t.Error(err)
	}

}

func TestSortTotal(t *testing.T) {
	var testmap map[string]int
	testmap = make(map[string]int)
	testmap["abcdefg"] = 3
	testmap["hijklmn"] = 5
	testmap["opqrstu"] = 8
	testmap["vwxyz"] = 12
	_, err := SortTotal(testmap)
	if err != nil {
		t.Error(err)
	}

}
func TestTop3(t *testing.T) {
	testarray := make([]CountArray, 4)
	testarray[0].id = "abc##"
	testarray[0].total = 5
	testarray[1].id = "def$$"
	testarray[1].total = 7
	testarray[2].id = "ghi%%"
	testarray[2].total = 9
	testarray[3].id = "jkl&&"
	testarray[3].total = 3
	_, _, _, err := Top3(testarray)
	if err != nil {
		t.Error(err)
	}

}
func TestFindIDInfo(t *testing.T) {
	var teststring1, teststring2, teststring3 string
	teststring1 = "ghi%%"
	teststring2 = "def$$"
	teststring3 = "abc##"

	testmodel := make([]datamodel.Coffee, 4)
	testreview0 := make([]datamodel.Review, 2)
	testmodel[0].Id = "23"
	testmodel[0].Name = "goodcoffee"
	testmodel[0].Rate = 8.3
	testreview0[0].StoreId = "def$$"
	testreview0[0].Text = "內裝舒適，座位寛敞，聊天小歇的好地方。"
	testmodel[0].Reviews = testreview0

	testreview1 := make([]datamodel.Review, 2)
	testmodel[1].Id = "99"
	testmodel[1].Name = "hellocoffee"
	testmodel[1].Rate = 9.9
	testreview1[0].StoreId = "abc##"
	testreview1[0].Text = "內裝舒適，座位寛敞，聊天小歇的好地方。"
	testmodel[1].Reviews = testreview1

	testreview2 := make([]datamodel.Review, 2)
	testmodel[2].Id = "45"
	testmodel[2].Name = "badcoffee"
	testmodel[2].Rate = 0.5
	testreview2[0].StoreId = "ghi%%"
	testreview2[0].Text = "內裝舒適，座位寛敞，聊天小歇的好地方。"
	testmodel[2].Reviews = testreview2

	_, _, _, err := FindIDInfo(teststring1, teststring2, teststring3, testmodel)
	if err != nil {
		t.Error(err)
	}

}
