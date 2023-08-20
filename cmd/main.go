package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/prondos/axdb/pkg/db"
	"github.com/prondos/axdb/pkg/filestorage"
	"github.com/prondos/axdb/pkg/rest"
)

type Data struct {
	Value string `json:"value" maxBytes:"1024"`
}

var (
	server  *rest.Server[string, filestorage.FileStorageMetadata, Data]
	service *rest.Service[string, filestorage.FileStorageMetadata, Data]
	table   *db.Table[string, filestorage.FileStorageMetadata, Data]
)

func init() {
	cwd, _ := os.Getwd()
	// Initialize the file storage for the application.
	storage := filestorage.NewFileStorage[string, filestorage.FileStorageMetadata, Data](filepath.Join(cwd, "storage"))

	// Create a new database table using the initialized file storage.
	table = db.NewTable[string, filestorage.FileStorageMetadata, Data](storage)

	// Create a new REST service using the created database table.
	service = rest.NewService[filestorage.FileStorageMetadata, Data](table)

	// Create a new REST server
	server = rest.NewServer[filestorage.FileStorageMetadata, Data](service)
}

func main() {
	// Open the database table.
	table.Open()
	defer table.Close()

	// start the server
	err := server.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
