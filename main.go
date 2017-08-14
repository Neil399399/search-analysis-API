package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"search-analysis-API/datamodel"
	"strconv"
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
	port   = "80"
	Search datamodel.Search
)

func main() {
	APIKey = "AIzaSyCigqPQLr341O-UL_jyJQNdX76fO0TtywA"
	keyword = "coffee"

	//wwww.google.com/maps?long=30&lat=30
	//http server
	myFunction := func() {
		//handle
		//&LAT=%f&LNG=%f&KEYWORD=%S", APIKey, Lat, Lng, keyword,
		http.HandleFunc("/search", DataSearch)
		//http.HandleFunc("/analysis", DataAnalysis)
		http.HandleFunc("/search-analysis", DataSearch_Analysis)

		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			panic("Connect Fail:" + err.Error())
		}
	}
	go myFunction()
	// use go channel to continous code
	endChannel := make(chan os.Signal)
	signal.Notify(endChannel)
	sig := <-endChannel
	fmt.Println("END!:", sig)

	//handle
}

/*
1 search: /search?lat... <- list results
2 analyze:/search?listresults <- get analysis
3 create UI


*/

func DataSearch(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	lat := req.FormValue("LAT")
	lng := req.FormValue("LNG")
	keyword := req.FormValue("KEYWORD")

	//check lat,lng format from http
	if len(lat) != 0 {
		lat64, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			fmt.Println("LAT has wrong format !!!")
			return
		}
		Search.LAT = lat64
	}
	if len(lng) != 0 {
		lng64, err := strconv.ParseFloat(lng, 64)
		if err != nil {
			fmt.Println("LNG has wrong format !!!")
			return
		}
		Search.LNG = lng64
	}

	//check from client
	Search.APIKEY = APIKey
	Search.KEYWORD = keyword
	if !Search.Verify(Search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "json")

	//search
	List, err := PlaceSearch(keyword, Search.LAT, Search.LNG)
	if err != nil {
		fmt.Println("google Place Search Error!!", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Println("google Search Success!!")
	//convert to json, give to fprint
	b, err := json.Marshal(List)
	if err != nil {
		fmt.Println("Json Marchal Error!!", err)
	}

	fmt.Fprint(w, string(b))
}

/*
func DataAnalysis(w http.ResponseWriter, req *http.Request) {
	var top [3]string
	var name1, name2, name3 string
	querys = []string{name1, name2, name3}
	if req.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")
		name1 = req.FormValue("analysis word 1")
		name2 = req.FormValue("analysis word 2")
		name3 = req.FormValue("analysis word 3")

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		type RequestMessage struct {
			Params []string
			Data   []datamodel.Coffee
		}

		var requestMessage RequestMessage
		//check err
		err = json.Unmarshal(body, &requestMessage)
		if err != nil {
			fmt.Println("Json Unmarshal Error!!", err)
			return
		}
		//check err

		//run jieba
		jiebres, err := jiebatest(list, querys)
		if err != nil {
			fmt.Println("jieba Error!!", err)
		}
		//count total
		sortres, err := SortTotal(jiebres)
		if err != nil {
			fmt.Println("Sort Total Error!!", err)
		}

		//find top3
		first, second, third, err := Top3(sortres)
		if err != nil {
			fmt.Println("Find Top3 Error!!", err)
		}

		//print top3

		top1, top2, top3, err := FindIDInfo(first, second, third, list)
		if err != nil {
			fmt.Println("json marshal failed!!", err)
		}
		top[0] = top1
		top[1] = top2
		top[2] = top3
		b, err := json.Marshal(top)
		fmt.Fprint(w, string(b))
		if err != nil {
			fmt.Println("Find ID Info Error!!", err)
		}

	}
}
*/
func DataSearch_Analysis(w http.ResponseWriter, req *http.Request) {
	var top [3]string
	//set header to tell server which http domain can connect
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	lat := req.FormValue("LAT")
	lng := req.FormValue("LNG")
	keyword := req.FormValue("KEYWORD")
	name1 := req.FormValue("analysis_word1")
	name2 := req.FormValue("analysis_word2")
	name3 := req.FormValue("analysis_word3")
	querys = []string{name1, name2, name3}

	//check lat,lng format from http
	if len(lat) != 0 {
		lat64, err := strconv.ParseFloat(lat, 64)
		if err != nil {
			fmt.Println("LAT has wrong format !!!")
			return
		}
		Search.LAT = lat64
	}
	if len(lng) != 0 {
		lng64, err := strconv.ParseFloat(lng, 64)
		if err != nil {
			fmt.Println("LNG has wrong format !!!")
			return
		}
		Search.LNG = lng64
	}

	//check from client
	Search.APIKEY = APIKey
	Search.KEYWORD = keyword
	if !Search.Verify(Search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "json")
	//search
	List, err := PlaceSearch(keyword, Search.LAT, Search.LNG)
	if err != nil {
		fmt.Println("google Place Search Error!!", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	fmt.Println("google  Search Success!!")
	//Analysis
	jiebres, err := jiebatest(List, querys)
	if err != nil {
		fmt.Println("jieba Error!!", err)
	}
	//count total
	sortres, err := SortTotal(jiebres)
	if err != nil {
		fmt.Println("Sort Total Error!!", err)
	}

	//find top3
	first, second, third, err := Top3(sortres)
	if err != nil {
		fmt.Println("Find Top3 Error!!", err)
	}

	//print top3

	top1, top2, top3, err := FindIDInfo(first, second, third, List)
	if err != nil {
		fmt.Println("json marshal failed!!", err)
	}
	top[0] = top1
	top[1] = top2
	top[2] = top3
	b, err := json.Marshal(top)
	fmt.Fprint(w, string(b))
	if err != nil {
		fmt.Println("Find ID Info Error!!", err)
	}

}
