package main

import (
	"database/sql"
	"fmt"
	"net"
	"search-analysis-API/datamodel"
)

type Storage interface {
	Read(id string) (datamodel.Coffee, error)
	Write(data datamodel.Coffee) error
	ReadId(data datamodel.Coffee) (*sql.Rows, error)
	ReadName(data datamodel.Coffee) (*sql.Rows, error)
	ReadPlaceID(data datamodel.Coffee) (*sql.Rows, error)
	//	ReadReviewsByID(id string) []string
}

type StoppableListener struct {
	*net.TCPListener          //Wrapped listener
	stop             chan int //Channel used only to indicate listener should shutdown
}
type StoreImpl struct {
}

func (s StoreImpl) Read(id string) (datamodel.Coffee, error) {
	fmt.Println("READ?")

	return datamodel.Coffee{}, nil
}

func (s StoreImpl) Write(data datamodel.Coffee) error {
	fmt.Println("WRITE?")

	return nil
}

//[]string
func (s StoreImpl) ReadId(data datamodel.Coffee) (*sql.Rows, error) {
	fmt.Println("ReadId?")

	return nil, nil
}

func (s StoreImpl) ReadName(data datamodel.Coffee) (*sql.Rows, error) {
	fmt.Println("ReadName?")
	return nil, nil
}

func (s StoreImpl) ReadPlaceID(data datamodel.Coffee) (*sql.Rows, error) {
	fmt.Println("ReadPlaceID?")
	return nil, nil
}
