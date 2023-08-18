package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prondos/axdb/pkg/db"
	"github.com/prondos/axdb/pkg/filestorage"
	"github.com/prondos/axdb/pkg/rest"
)

type Data struct {
	Value string `json:"value" maxBytes:"1024"`
}

func main() {
	storage := filestorage.NewFileStorage[string, filestorage.FileStorageMetadata, Data]("../storage")
	table := db.NewTable[string, filestorage.FileStorageMetadata, Data](storage)
	service := rest.NewService[filestorage.FileStorageMetadata, Data](table)
	table.Open()
	defer table.Close()

	data1 := &Data{Value: "Data 1 value"}
	data2 := &Data{Value: "Data 2 value"}
	if err := table.Insert("key1", *data1); err != nil {
		log.Printf("error inserting %v, %v", *data1, err)
	}
	if err := table.Insert("key2", *data2); err != nil {
		log.Printf("error inserting %v, %v", *data2, err)
	}

	router := gin.Default()
	router.GET("/items", service.Index)
	router.GET("/items/:key", service.Get)
	router.PUT("/items/:key", service.Put)
	router.Run("localhost:6600")
}
