package sqlstore

import (
	"database/sql"
	"fmt"
	"search-analysis-API/datamodel"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

//"root:123456@tcp(localhost:3306)/hello"
var (
	cof datamodel.Coffee
)

type WriteToSQL struct {
	serverurl string
	database  string
	db        *sql.DB
}

func NewWriteToSQL(username, password, serverurl, database string) *WriteToSQL {
	// Create the database handle, confirm driver is present
	DB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4,utf8", username, password, serverurl, database))
	if err != nil {
		fmt.Println("Connect Error!!", err)
		panic(err)
	}

	return &WriteToSQL{db: DB}
}

func (w *WriteToSQL) Read(id string) (datamodel.Coffee, error) {
	return datamodel.Coffee{}, nil
}

func (w *WriteToSQL) ReadId(data datamodel.Coffee) (*sql.Rows, error) {
	return w.read("SELECT Comment FROM CoffeeComment WHERE PlaceID IN (SELECT PlaceID FROM CoffeeInfo WHERE Rate= CAST(? AS DECIMAL))", data.Rate)
}

func (w *WriteToSQL) ReadName(data datamodel.Coffee) (*sql.Rows, error) {

	return w.read("SELECT Name FROM CoffeeInfo WHERE PlaceID = ?", data.Name)
}

func (w *WriteToSQL) ReadPlaceID(data datamodel.Coffee) (*sql.Rows, error) {
	int_ID, err := strconv.Atoi(data.Id)
	if err != nil {
		fmt.Println("ReadPlaceID String to Int Error!!", err)
	}
	return w.read("SELECT PlaceID FROM CoffeeComment WHERE ID=?", int_ID)
}

func (w *WriteToSQL) read(sqlQuery string, args ...interface{}) (*sql.Rows, error) {
	res, err := w.db.Query(sqlQuery, args...)
	if err != nil {
		fmt.Println("Read Comment Error!!", err)
	}

	//make a slice of the data
	//[]string{}
	for res.Next() {
		tmp := ""
		err := res.Scan(&tmp)
		if err != nil {
			fmt.Println("SQL Result Print Error!!", err)
		}
		fmt.Println(tmp)
		//slice = append(slice, tmp)
		//append to []string{}
	}

	return res, nil
}

func (w *WriteToSQL) Write(data datamodel.Coffee) error {

	_, err := w.db.Exec("INSERT INTO CoffeeInfo (PlaceID,Name,Rate) VALUES (?,?,?)", data.Id, data.Name, data.Rate)
	if err != nil {
		fmt.Println("Write Info Error!!")
		panic(err)
	}

	for i := 1; i < len(data.Reviews); i++ {

		_, err = w.db.Exec("INSERT INTO CoffeeComment (PlaceID,Comment) VALUES (?,?)", data.Reviews[i].StoreId, data.Reviews[i].Text)
		if err != nil {
			fmt.Println("Write Comment Error!!")
			fmt.Println("Index: ", i)
			fmt.Println(data.Reviews[i].Text)
			panic(err)
		}

	}

	return nil
}

/*
func (w *WriteToSQL) ReadReviewsByID(id string) []string {

}
*/
