package filestorage

import (
	"fmt"
	"io"
	"log"

	"github.com/prondos/axdb/pkg/db"
)

func (p *FileStorage[IT, DT]) LoadAll() []*db.Record[IT, DT] {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	var records []*db.Record[IT, DT]

	for {
		index, offset, err := p.indexFromReader(p.IndexReader)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(fmt.Sprintf("error reading record index: %v", err))
		}
		p.Index[index] = &FileStorageMetadata{
			stored: true,
			offset: offset,
		}
		record := &db.Record[IT, DT]{
			Index: index,
		}
		record.Data, err = p.dataFromReader(p.DataReader, offset)
		if err != nil {
			panic(fmt.Sprintf("error reading record data: %v", err))
		}

		records = append(records, record)
	}

	return records
}

func (p *FileStorage[IT, DT]) Store(record *db.Record[IT, DT]) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.Index[record.Index]; !ok {
		p.Index[record.Index] = &FileStorageMetadata{stored: false}
	}
	p.storeWaitGroup.Add(1)
	p.storeChannel <- record
	return nil
}

func (p *FileStorage[IT, DT]) Delete(record db.Record[IT, DT]) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return nil
}

func (p *FileStorage[IT, DT]) Close() {
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
