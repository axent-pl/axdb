package db

import (
	"fmt"
	"sync"
)

type Table[IT comparable, DT any] struct {
	mutex   sync.Mutex
	storage StorageInterface[IT, DT]
	records map[IT]*Record[IT, DT]
}

func NewTable[IT comparable, DT any](storage StorageInterface[IT, DT]) *Table[IT, DT] {
	tab := &Table[IT, DT]{
		storage: storage,
		records: make(map[IT]*Record[IT, DT]),
	}
	return tab
}

func (tab *Table[IT, DT]) List() []IT {
	var ret []IT
	for index := range tab.records {
		ret = append(ret, index)
	}
	return ret
}

func (tab *Table[IT, DT]) Read(index IT) (DT, error) {
	if rec, ok := tab.records[index]; ok {
		return rec.Data, nil
	} else {
		return *new(DT), fmt.Errorf("record with index %v not found", index)
	}
}

func (tab *Table[IT, DT]) Insert(index IT, data DT) error {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	if _, ok := tab.records[index]; ok {
		return fmt.Errorf("record with index %v already exists", index)
	}
	rec := NewRecord[IT, DT](index, data)
	tab.records[index] = rec
	tab.storage.Store(rec)

	return nil
}

func (tab *Table[IT, DT]) Update(index IT, data DT) error {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	if rec, ok := tab.records[index]; ok {
		rec.Data = data
		tab.storage.Store(rec)
	}

	return fmt.Errorf("record with index %v does not exist", index)
}

func (tab *Table[IT, DT]) InsertOrUpdate(index IT, data DT) error {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	if rec, ok := tab.records[index]; ok {
		rec.Data = data
		tab.storage.Store(rec)
	} else {
		rec := NewRecord[IT, DT](index, data)
		tab.records[index] = rec
		tab.storage.Store(rec)
	}
	return nil
}

func (tab *Table[IT, DT]) Delete(index IT) error {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	if rec, ok := tab.records[index]; ok {
		if err := tab.storage.Delete(*rec); err != nil {
			return err
		}
		delete(tab.records, index)
	}

	return fmt.Errorf("record with index %v does not exist", index)
}

func (tab *Table[IT, DT]) Open() {
	for _, record := range tab.storage.LoadAll() {
		tab.records[record.Index] = record
	}

}

func (tab *Table[IT, DT]) Close() {
	tab.storage.Close()
}
