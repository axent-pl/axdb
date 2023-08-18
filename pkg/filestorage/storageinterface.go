package filestorage

import (
	"fmt"
	"io"
	"log"

	"github.com/prondos/axdb/pkg/db"
)

func (p *FileStorage[IT, MT, DT]) LoadAll() []*db.Record[IT, FileStorageMetadata, DT] {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var records []*db.Record[IT, FileStorageMetadata, DT]

	for {
		index, offset, err := p.indexFromReader(p.IndexReader)
		record := &db.Record[IT, FileStorageMetadata, DT]{
			Index: index,
			Metadata: &FileStorageMetadata{
				offset: offset,
			},
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Sprintf("error reading record index: %v", err))
		}
		record.Data, err = p.dataFromReader(p.DataReader, offset)
		if err != nil {
			panic(fmt.Sprintf("error reading record data: %v", err))
		}

		records = append(records, record)
	}

	return records
}

func (p *FileStorage[IT, MT, DT]) Store(record *db.Record[IT, FileStorageMetadata, DT]) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if record.Metadata == nil {
		record.Metadata = &FileStorageMetadata{stored: false}
	}
	p.storeWaitGroup.Add(1)
	p.storeChannel <- record
	return nil
}

func (p *FileStorage[IT, MT, DT]) Delete(record db.Record[IT, FileStorageMetadata, DT]) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return nil
}

func (p *FileStorage[IT, MT, DT]) Close() {
	log.Print("Closing file storage")
	p.storeWaitGroup.Wait()
	if err := p.DataWriter.Close(); err != nil {
		log.Printf("Error closing file storage DataWriter: %v", err)
	}
	if err := p.DataReader.Close(); err != nil {
		log.Printf("Error closing file storage DataReader: %v", err)
	}
	if err := p.IndexWriter.Close(); err != nil {
		log.Printf("Error closing file storage IndexWriter: %v", err)
	}
	if err := p.IndexReader.Close(); err != nil {
		log.Printf("Error closing file storage IndexReader: %v", err)
	}
	log.Print("Closing file storage DONE")
}
