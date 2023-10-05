package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/prondos/axdb/pkg/db"
	"github.com/prondos/axdb/pkg/filestorage"
	"github.com/prondos/axdb/pkg/rest"
)

type Data struct {
	Value string `json:"value" maxBytes:"1024"`
}

var (
	signalChannel chan os.Signal
	server        *rest.Server[string, Data]
	service       *rest.Service[string, Data]
	table         *db.Table[string, Data]
)

func init() {
	cwd, _ := os.Getwd()
	// Initialize the file storage for the application.
	storage := filestorage.MustNewFileStorage[string, Data](filestorage.WithDatadir(filepath.Join(cwd, "storage")))

	// Create a new database table using the initialized file storage.
	table = db.NewTable[string, Data](storage)

	// Create a new REST service using the created database table.
	service = rest.NewService[Data](table)

	// Create a new REST server
	server = rest.NewServer[Data](service)

	// Configure signalChannel
	signalChannel = make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
}

func main() {
	// Open the database table.
	if err := table.Open(); err != nil {
		os.Exit(1)
	}
	defer table.Close()

	// Init cancellation context triggered by SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-signalChannel
		cancel()
	}()
	defer cancel()

	// start the server
	if err := server.Start(ctx); err != nil {
		os.Exit(1)
	}
}
