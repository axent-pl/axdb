package db

import (
	"fmt"
	"sync"
)

type Table[IT comparable, MT any, DT any] struct {
	mutex   sync.Mutex
	storage StorageInterface[IT, MT, DT]
	records map[IT]*Record[IT, MT, DT]
}

func NewTable[IT comparable, MT any, DT any](storage StorageInterface[IT, MT, DT]) *Table[IT, MT, DT] {
	tab := &Table[IT, MT, DT]{
		storage: storage,
		records: make(map[IT]*Record[IT, MT, DT]),
	}
	return tab
}

func (tab *Table[IT, MT, DT]) List() []IT {
	var ret []IT
	for index := range tab.records {
		ret = append(ret, index)
	}
	return ret
}

func (tab *Table[IT, MT, DT]) Read(index IT) (DT, error) {
	if rec, ok := tab.records[index]; ok {
		return rec.Data, nil
	} else {
		return *new(DT), fmt.Errorf("record with index %v not found", index)
	}
}

func (tab *Table[IT, MT, DT]) Insert(index IT, data DT) error {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	if _, ok := tab.records[index]; ok {
		return fmt.Errorf("record with index %v already exists", index)
	}
	rec := NewRecord[IT, MT, DT](index, data)
	tab.records[index] = rec
	tab.storage.Store(rec)

	return nil
}

func (tab *Table[IT, MT, DT]) Update(index IT, data DT) error {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	if rec, ok := tab.records[index]; ok {
		rec.Data = data
		tab.storage.Store(rec)
	}

	return fmt.Errorf("record with index %v does not exist", index)
}

func (tab *Table[IT, MT, DT]) InsertOrUpdate(index IT, data DT) error {
	tab.mutex.Lock()
	defer tab.mutex.Unlock()

	if rec, ok := tab.records[index]; ok {
		rec.Data = data
		tab.storage.Store(rec)
	} else {
		rec := NewRecord[IT, MT, DT](index, data)
		tab.records[index] = rec
		tab.storage.Store(rec)
	}
	return nil
}

func (tab *Table[IT, MT, DT]) Delete(index IT) error {
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

func (tab *Table[IT, MT, DT]) Open() {
	for _, record := range tab.storage.LoadAll() {
		tab.records[record.Index] = record
	}

}

func (tab *Table[IT, MT, DT]) Close() {
	tab.storage.Close()
}
