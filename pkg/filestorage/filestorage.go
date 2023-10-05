package filestorage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/prondos/axdb/pkg/db"
	"github.com/prondos/axdb/pkg/utils"
)

const (
	datafile string = "data.kvp"
	indxfile string = "indx.kvp"
)

type IndexWrapper[IT comparable] struct {
	I IT
	O int64
}

type DataWrapper[DT any] struct {
	D DT
}

type FileStorageMetadata struct {
	stored bool
	offset int64
}

type FileStorage[IT comparable, DT any] struct {
	mutex          sync.Mutex
	offset         int64
	dataBlockSize  int
	Datadir        string
	DataWriter     *os.File
	DataReader     *os.File
	IndexWriter    *os.File
	IndexReader    *os.File
	Index          map[IT]*FileStorageMetadata
	storeWaitGroup sync.WaitGroup
	storeChannel   chan *db.Record[IT, DT]
	deleteChannel  chan int64
}

type FileStorageConfig struct {
	datadir string
}

func NewFileStorageConfig() *FileStorageConfig {
	return &FileStorageConfig{
		datadir: "storage",
	}
}

func WithDatadir(datadir string) func(*FileStorageConfig) {
	return func(c *FileStorageConfig) {
		c.datadir = datadir
	}
}

func MustNewFileStorage[IT comparable, DT any](options ...func(*FileStorageConfig)) *FileStorage[IT, DT] {
	calculatedOptions := NewFileStorageConfig()
	for _, option := range options {
		option(calculatedOptions)
	}

	p := &FileStorage[IT, DT]{
		Datadir:       calculatedOptions.datadir,
		Index:         make(map[IT]*FileStorageMetadata),
		storeChannel:  make(chan *db.Record[IT, DT], 1),
		deleteChannel: make(chan int64, 1),
	}

	if err := p.init(); err != nil {
		panic(err)
	}

	return p
}

func (p *FileStorage[IT, DT]) init() error {
	var err error

	// calculate dataa block size
	p.dataBlockSize, err = utils.GetBlockSize[DT]()
	if err != nil {
		return err
	}

	// check datadir
	datadirStat, err := os.Stat(p.Datadir)
	if err != nil {
		return err
	}
	if !datadirStat.IsDir() {
		return fmt.Errorf("%v is not a directory", p.Datadir)
	}

	// init data file writer
	dataPath := filepath.Join(p.Datadir, datafile)
	p.DataWriter, err = os.OpenFile(dataPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	// init data file reader
	p.DataReader, err = os.Open(dataPath)
	if err != nil {
		return err
	}
	p.offset = p.dataWriterOffset()

	// init index file writer
	indexPath := filepath.Join(p.Datadir, indxfile)
	p.IndexWriter, err = os.OpenFile(indexPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	// init index file reader
	p.IndexReader, err = os.Open(indexPath)
	if err != nil {
		return err
	}

	go p.processStoreChannel()

	return nil
}

func (p *FileStorage[IT, DT]) dataWriterOffset() int64 {
	writerStat, _ := p.DataWriter.Stat()
	return writerStat.Size()
}

func (p *FileStorage[IT, DT]) indexWriterOffset() int64 {
	writerStat, _ := p.IndexWriter.Stat()
	return writerStat.Size()
}

func (p *FileStorage[IT, DT]) processStoreChannel() {
	for {
		r := <-p.storeChannel

		if !p.Index[r.Index].stored {
			p.Index[r.Index].offset = p.offset
			indexBytes := p.indexToBytes(r.Index, p.Index[r.Index].offset)
			p.IndexWriter.Write(indexBytes)
		}

		dataBytes := p.dataToBytes(r.Data, p.dataBlockSize)
		_, err := p.DataWriter.WriteAt(dataBytes, p.Index[r.Index].offset)
		if err != nil {
			panic(err)
		}
		p.offset = p.dataWriterOffset()
		p.storeWaitGroup.Done()
	}
}
