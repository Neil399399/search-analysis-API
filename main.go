package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"search-analysis-API/datamodel"
	"search-analysis-API/search"
	"strconv"
)

var (
	port   = "80"
	Search datamodel.Search
)

func main() {

	//http server
	myFunction := func() {
		//handle
		http.HandleFunc("/search", DataSearch)
		http.HandleFunc("/analysis", DataAnalysis)
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
}

//function
func DataSearch(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	apikey := req.FormValue("APIKEY")
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
	err := search.Initialize(apikey, 500)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	search500 := search.NewSearch(apikey, 500)

	Search.KEYWORD = keyword
	if !Search.Verify(Search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "json")

	//search
	List, err := search500.Place(keyword, Search.LAT, Search.LNG)
	if err != nil {
		fmt.Println("google Place Search Error!!", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//fmt.Println("google Search Success!!")

	//convert to json, give to fprint
	b, err := json.Marshal(List)
	if err != nil {
		fmt.Println("Json Marchal Error!!", err)
	}
	fmt.Fprint(w, string(b))
}

func DataAnalysis(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	var top [3]string

	if req.Method == "POST" {
		w.Header().Set("Content-Type", "application/json")

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
		//run jieba
		jiebres, err := jiebatest(requestMessage.Data, requestMessage.Params)
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
		top1, top2, top3, err := FindIDInfo(first, second, third, requestMessage.Data)
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

func DataSearch_Analysis(w http.ResponseWriter, req *http.Request) {
	var top [3]string
	//set header to tell server which http domain can connect
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	apikey := req.FormValue("APIKEY")
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
	err := search.Initialize(apikey, 500)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	search500 := search.NewSearch(apikey, 500)
	Search.KEYWORD = keyword
	if !Search.Verify(Search) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	//header set
	w.Header().Set("Content-Type", "json")
	//search
	List, err := search500.Place(keyword, Search.LAT, Search.LNG)
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
