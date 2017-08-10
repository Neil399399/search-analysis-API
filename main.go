package main

import (
	"fmt"
)

/*
	create blevedir+"/"+lat_long_timestamp
	check if used
	if not, create index
		input lat/long data into bleve index
	do jieba
	return results
*/
var (
	index_dir       string
	input_indexname string
)

func main() {
	//search
	coffeeList, err := PlaceSearch()
	if err != nil {
		fmt.Println("google Place Search Error!!", err)
	}
	//filterResults(coffeeList)  //filter.go
	/*
		var data []datamodel.Comment
		n := 0
		for i := 0; i < len(coffeeList); i++ {
			for j := 0; j < len(coffeeList[i].Reviews); j++ {

				data[n+j].PlaceID = coffeeList[i].Reviews[j].StoreId
			}
			//fmt.Println(coffeeList[cof].Reviews)
			n = n + 5
	*/

	//create index
	index_dir = "random"
	/*
		err = CreateIndex(coffeeList, index_dir)
		if err != nil {
			fmt.Println("CreateIndex Error!!", err)
		}
	*/
	//run jieba
	jieb_res, err := jiebatest(index_dir, coffeeList)
	if err != nil {
		fmt.Println("jieba Error!!", err)
	}
	err = CountResult(jieb_res)

	/*
		//count total
		sort_res, err := SortTotal(jieb_res)
		if err != nil {
			fmt.Println("Sort Total Error!!", err)
		}
		//find top3
		first, second, third, err := Top3(sort_res)
		if err != nil {
			fmt.Println("Find Top3 Error!!", err)
		}
		//print top3
		err = FindIDInfo(first)
		if err != nil {
			fmt.Println("Find ID Info Error!!", err)
		}
		err = FindIDInfo(second)
		if err != nil {
			fmt.Println("Find ID Info Error!!", err)
		}
		err = FindIDInfo(third)
		if err != nil {
			fmt.Println("Find ID Info Error!!", err)
		}
	*/
}