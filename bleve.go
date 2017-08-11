package main

import (
	"encoding/json"
	"fmt"
	"search-analysis-API/datamodel"
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
	querys []string
	index  bleve.Index
)

func init() {
	//bleve.Open
	indexMapping := bleve.NewIndexMapping()
	_, err := bleve.NewMemOnly(indexMapping)

	if err != nil {
		panic("bleve open failed!!" + err.Error())
	}
	err = indexMapping.AddCustomTokenizer("gojieba",
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
		panic("bleve open failed!!" + err.Error())
	}

}

func getFreeIndex() bleve.Index {

	return index
}

func jiebatest(com []datamodel.Coffee, querys []string) (map[string]int, error) {
	type Result struct {
		Id    string
		Score float64
	}

	//create index_dir
	for i := 0; i < len(com); i++ {
		err := getFreeIndex().Index(com[i].Id, com[i].Reviews)
		if err != nil {

		}
	}

	dataCounter := make(map[string]int)
	for _, q := range querys {

		req := bleve.NewSearchRequest(bleve.NewQueryStringQuery(q))

		req.Highlight = bleve.NewHighlight()
		res, err := getFreeIndex().Search(req)
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
	getFreeIndex().Close()
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

	var top1, top2, top3 string
	top1 = array[0].id
	top2 = array[1].id
	top3 = array[2].id

	return top1, top2, top3, nil
}

func FindIDInfo(first, second, third string, com []datamodel.Coffee) error {

	for idx, cof := range com {
		if len(cof.Reviews) > 0 {
			if first == cof.Reviews[0].StoreId {
				fmt.Println("TOP1", com[idx].Name)
			}
			if second == cof.Reviews[0].StoreId {
				fmt.Println("TOP2", com[idx].Name)
			}
			if third == cof.Reviews[0].StoreId {
				fmt.Println("TOP3", com[idx].Name)
			}
		}
	}
	return nil
}
