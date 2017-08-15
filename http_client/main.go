package main

import (
	"fmt"
	"net/http"
)

type Client struct {
	url         string
	contentType string
}

var (
	url           = "http://localhost"
	ngrokFrontend = "http://67ddbfbe.ngrok.io"
	ngrokBackend  = " http://22f820fd.ngrok.io"
	port          = "8080"
)

func init() {}

func main() {
	//handle
	//&LAT=%f&LNG=%f&KEYWORD=%S", APIKey, Lat, Lng, keyword,
	http.HandleFunc("/form", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, WebView)
	})

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic("Connect Fail:" + err.Error())
	}
	/*
		client := &Client{url: url, contentType: "json"}
		err = client.GetSearchandAnalaysis()
	*/
}

/*
func GetFromCmd(client *Client) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println(" Use or Quit: ")
		input, _, _ := reader.ReadLine()
		inputStr := string(input)
		if inputStr == "Quit" {
			fmt.Println("Quit  System !!")
			break
		} else if inputStr == "Use" {
			//search ,analysis or search-analysis
			fmt.Println("Search ,Analysis or Search-Analysis :")
			q, _, _ := reader.ReadLine()
			str := string(q)
			if str == "Search" {
				client.GetSearch()
			} else if str == "Analysis" {

			} else if str == " Search-Analysis" {

			}

		}

	}

}
*/
/*
func (c *Client) GetSearchandAnalaysis() error {
	request, err := http.NewRequest("GET", c.url+"/search-analysis", nil)

	return nil
}*/

/*
func (c *Client) GetSearch() error {

	reader := bufio.NewReader(os.Stdin)
	request, err := http.NewRequest("GET", c.url+"/search", nil)
	fmt.Println("Input your Keyword:")
	keyword, _, _ := reader.ReadLine()
	keywordstr := string(keyword)

	fmt.Println("Input your loaction(lat and lng)")
	fmt.Println("lat:")
	lat, _, _ := reader.ReadLine()
	latstr := string(lat)
	if err != nil {
		fmt.Println("Request Error!!")
		panic(err.Error())
	}
	fmt.Println("lng:")
	lng, _, _ := reader.ReadLine()
	lngstr := string(lng)
	if err != nil {
		fmt.Println("Request Error!!")
		panic(err.Error())
	}

	queryForm := request.URL.Query()
	queryForm.Add("KEYWORD", keywordstr)
	queryForm.Add("LAT", latstr)
	queryForm.Add("LNG", lngstr)
	request.URL.RawQuery = queryForm.Encode()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("connect error!!: ", err)
	}
	if response.StatusCode != http.StatusOK {
		//decode body into []datamodel.Coffee
		//type RequestMessage struct {
		//
		//}

		return errors.New(response.Status)
	}
	return nil
}
*/
