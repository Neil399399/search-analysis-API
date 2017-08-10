package main

import (
	"encoding/json"
	"fmt"
	"search-analysis-API/datamodel"
	"search-analysis-API/storage/sqlstore"
	"sort"

	"github.com/blevesearch/bleve"
	"github.com/yanyiwu/gojieba"
)

type JiebaTokenizer struct {
	handle *gojieba.Jieba
}

type CountArray struct {
	id    string
	total int
}

var (
	name    string
	placeID string
)

func jiebatest(index_dir string, com []datamodel.Coffee) (map[string]int, error) {
	type Result struct {
		Id    string
		Score float64
	}

	indexMapping := bleve.NewIndexMapping()
	err := indexMapping.AddCustomTokenizer("gojieba",
		map[string]interface{}{
			"dictpath":   "jieba/dict.txt",
			"hmmpath":    "jieba/hmm_model.utf8",
			"idf":        "jieba/idf.utf8",
			"stop_words": "jieba/stop_word.utf8",
			"type":       "unicode",
		},
	)

	if err != nil {
		fmt.Println("Tokenizer Error!!", err)
	}
	/*
		err = indexMapping.AddCustomAnalyzer("gojieba",
			map[string]interface{}{
				"type":      "gojieba",
				"tokenizer": "gojieba",
			},
		)
		if err != nil {
			panic(err)
		}
	*/
	indexMapping.DefaultAnalyzer = "gojieba"

	querys := []string{
		"環境舒服",
		"服務好",
		"咖啡好喝",
		"好喝",
		"好",
	}

	//create index_dir
	index, err := bleve.New(index_dir, indexMapping)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(com); i++ {
		err = index.Index(com[i].Id, com[i].Reviews)
		if err != nil {

		}
	}

	/*
		index, err := bleve.Open(index_dir)
		if err != nil {
			fmt.Println("Open index Error!!", err)
		}
	*/
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
	index.Close()
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

	var top,top2,top3 string
	top1 = array[0].id
	top2 = array[1].id
	top3 = array[2].id

	}

	return top1, top2, top3, nil
}

func FindIDInfo(first string, com []datamodel.Coffee) error {

	for i:=0;i<len(com); i++{
		if first==com[i].Reviews[0].StoreId{
		fmt.Println(com[i].Name)
		}

	}
	return nil
}
