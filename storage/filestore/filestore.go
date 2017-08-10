package filestore

import (
	"encoding/json"
	"fmt"
	"search-analysis-API/datamodel"

	"io/ioutil"
	"os"
	"strings"
)

type WriteInFile struct {
	filename string
}

func NewWriteInFile(filename string) (WriteInFile, error) {
	return WriteInFile{filename: filename}, nil
}

func (w WriteInFile) Write(data datamodel.Coffee) error {
	file, err := os.OpenFile(w.filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Open File Error!:", err)
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json Marshal Error!:", err)
	}
	// byte change to string and +\n,then change to byte again.
	str := string(b) + "/n"
	fmt.Println(str)
	_, err = file.WriteString(str)

	file.Close()
	return nil
}

func (w WriteInFile) Read() ([]datamodel.Coffee, error) {

	b, err := ioutil.ReadFile(w.filename)
	if err != nil {
		fmt.Println("read error: ", err)
	}
	// Change byte(b) to String(b) and find /n split.
	jsonArr := strings.Split(string(b), "/n")
	//Create new List and append
	coffees := []datamodel.Coffee{}
	// unmarshal each list
	for i := 0; i < len(jsonArr)-1; i++ {
		var data datamodel.Coffee
		//unmarshal and Change String to byte
		err = json.Unmarshal([]byte(jsonArr[i]), &data)
		if err != nil {
			fmt.Println("json err:", err)

		}

		coffees = append(coffees, data)

	}

	return coffees, nil
}
