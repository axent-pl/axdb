package main

import (
	"fmt"
	"log"

	"github.com/prondos/axdb/pkg/filestorage"
)

type Data struct {
	Name    string `json:"name" maxBytes:"100"`
	Comment string `json:"comment" maxBytes:"100"`
}

func main() {
	table := filestorage.NewTable[string, Data]("../storage")
	defer table.Close()

	data1 := &Data{Name: "John", Comment: "Nice"}
	data2 := &Data{Name: "John", Comment: "Nice2333ssssssss"}
	if err := table.Insert("key1", *data1); err != nil {
		log.Printf("error inserting %v, %v", *data1, err)
	}
	if err := table.Insert("key2", *data2); err != nil {
		log.Printf("error inserting %v, %v", *data2, err)
	}

	indices := table.List()
	for _, index := range indices {
		rec, _ := table.Read(index)
		fmt.Printf("record[%v] := %v\n", index, rec)
	}

}
