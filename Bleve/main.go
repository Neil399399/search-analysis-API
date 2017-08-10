package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"search-analysis-API/datamodel"
	"search-analysis-API/storage/sqlstore"
	"sort"

	"github.com/blevesearch/bleve"
)

/*
type JiebaTokenizer struct {
	handle *gojieba.Jieba
}
*/
type CountArray struct {
	id    string
	total int
}

var (
	filename   = "CoffeeComment.json"
	filename1  = "CoffeeComment.json"
	bleve_dir  = "bleve"
	index_dir  = "coffee.bleve"
	index_dir1 = "coffeeInfo.bleve"
	name       string
	placeID    string
)

func main() {
	/*
		create blevedir+"/"+lat_long_timestamp
		check if used
		if not, create index
			input lat/long data into bleve index
		do jieba
		return results
	*/

	/*
		com, err := Read(filename)
		if err != nil {
			fmt.Println("Read Error!!", err)
		}

		err = CreateIndex(com, index_dir1)
		if err != nil {
			fmt.Println("CreateIndex Error!!", err)
		}
	*/

	jieb_res, err := jiebatest(index_dir)
	if err != nil {
		fmt.Println("jieba Error!!", err)
	}
	sort_res, err := SortTotal(jieb_res)
	if err != nil {
		fmt.Println("Sort Total Error!!", err)
	}
	first, second, third, err := Top3(sort_res)
	if err != nil {
		fmt.Println("Find Top3 Error!!", err)
	}
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

}

func Read(filename string) ([]datamodel.Comment, error) {

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("read error: ", err)
	}
	//Create new List and append
	com := []datamodel.Comment{}
	// unmarshal each list

	//unmarshal and Change String to byte
	err = json.Unmarshal(b, &com)
	if err != nil {
		fmt.Println("json err:", err)
	}

	return com, nil
}

func CreateIndex(com []datamodel.Comment, index_dir string) error {
	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(index_dir, mapping)
	if err != nil {
		fmt.Println(err)
	}

	// index some data
	for i := 0; i < len(com); i++ {
		err = index.Index(com[i].ID, com[i].Comment)
		//fmt.Println(com[i].Comment)
	}

	return nil
}

func jiebatest(index_dir string) (map[string]int, error) {
	type Result struct {
		Id    string
		Score float64
	}
	indexMapping := bleve.NewIndexMapping()
	err := indexMapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dictpath":   "jieba/dict.txt",
			"hmmpath":    "jieba/hmm_model.utf8",
			"idf":        "idf.utf8",
			"stop_words": "stop_word.utf8",
			"type":       "unicode",
		},
	)
	if err != nil {
		fmt.Println("Tokenizer Error!!", err)
	}

	indexMapping.DefaultAnalyzer = "gojieba"

	querys := []string{
		"環境舒服",
		"不錯",
		"咖啡好喝",
		"好喝",
		"好",
	}

	index, err := bleve.Open(index_dir)
	if err != nil {
		fmt.Println("Open index Error!!", err)
	}
	dataCounter := make(map[string]int)
	for _, q := range querys {
		req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))
		req.Highlight = bleve.NewHighlight()
		res, err := index.Search(req)
		if err != nil {
			panic(err)
		}
		results := []Result{}
		for _, item := range res.Hits {
			results = append(results, Result{item.ID, item.Score})
		}

		for i := 0; i < len(results); i++ {
			dataCounter[results[i].Id]++
		}
	}

	return dataCounter, nil
}

func prettify(res *bleve.SearchResult) (string, error) {
	type Result struct {
		Id    string
		Score float64
	}
	results := []Result{}
	for _, item := range res.Hits {
		results = append(results, Result{item.ID, item.Score})
	}

	b, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	return string(b), nil
}

func CountResult(data map[string]int) error {

	for k, v := range data {
		fmt.Println("id:", k)
		fmt.Println("total:", v)
	}
	return nil
}

func SortTotal(data map[string]int) ([]CountArray, error) {

	//USE Slice (can't not use Array)
	countarrays := make([]CountArray, len(data))
	i := 0
	for k, v := range data {
		countarrays[i].id = k
		countarrays[i].total = v
		i++
		if i > len(data) {
			fmt.Println("i>len(data):", i)
		}
	}
	sort.Slice(countarrays, func(i, j int) bool {
		return countarrays[i].total >= countarrays[j].total
	})

	return countarrays, nil
}
func Top3(array []CountArray) (string, string, string, error) {
	sqlStore := sqlstore.NewWriteToSQL("root", "123456", "localhost", "hello")
	var top [3]string
	var Top map[string]interface{}
	Top = make(map[string]interface{})
	for i := 0; i < len(array); i++ {
		data := datamodel.Coffee{}
		data.Id = array[i].id
		sql_res, err := sqlStore.ReadPlaceID(data)
		if err != nil {
			fmt.Println("Search PlaceID in SQL Error!!", err)
		}
		for sql_res.Next() {
			err := sql_res.Scan(&placeID)
			if err != nil {
				fmt.Println("SQL Result Print Error!!", err)
			}
		}

		//find same placeID
		if Top[placeID] == nil {
			Top[placeID] = "excited"
		}
	}
	i := 0
	for k, _ := range Top {
		if i < 3 {
			top[i] = k
			i++
		} else {
			break
		}
	}

	return top[0], top[1], top[2], nil
}

func FindIDInfo(first string) error {

	data := datamodel.Coffee{}
	data.Name = first
	sqlStore := sqlstore.NewWriteToSQL("root", "123456", "localhost", "hello")
	sql_res, err := sqlStore.ReadName(data)
	if err != nil {
		fmt.Println("Search PlaceID in SQL Error!!", err)
	}
	for sql_res.Next() {
		err := sql_res.Scan(&name)
		if err != nil {
			fmt.Println("SQL Result Print Error!!", err)
		}

		fmt.Println(name)
	}
	return nil
}
