package db

import (
	"fmt"
	"sync"
)

type Table[IT comparable, DT any] struct {
	storage Storager[IT, DT]

	recordsMutex sync.Mutex
	records      map[IT]*Record[IT, DT]
}

func NewTable[IT comparable, DT any](storage Storager[IT, DT]) *Table[IT, DT] {
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
		return *new(DT), fmt.Errorf("%v %w", index, ErrNotFound)
	}
}

func (tab *Table[IT, DT]) Insert(index IT, data DT) error {
	tab.recordsMutex.Lock()
	defer tab.recordsMutex.Unlock()

	if _, ok := tab.records[index]; ok {
		return fmt.Errorf("%v %w", index, ErrExists)
	}
	rec := NewRecord[IT, DT](index, data)
	tab.records[index] = rec
	tab.storage.Store(rec)

	return nil
}

func (tab *Table[IT, DT]) Update(index IT, data DT) error {
	tab.recordsMutex.Lock()
	defer tab.recordsMutex.Unlock()

	if rec, ok := tab.records[index]; ok {
		rec.Data = data
		tab.storage.Store(rec)
	}

	return fmt.Errorf("%v %w", index, ErrNotFound)
}

func (tab *Table[IT, DT]) InsertOrUpdate(index IT, data DT) error {
	tab.recordsMutex.Lock()
	defer tab.recordsMutex.Unlock()

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
	tab.recordsMutex.Lock()
	defer tab.recordsMutex.Unlock()

	if rec, ok := tab.records[index]; ok {
		if err := tab.storage.Delete(*rec); err != nil {
			return err
		}
		delete(tab.records, index)
	}

	return fmt.Errorf("%v %w", index, ErrNotFound)
}

func (tab *Table[IT, DT]) Open() {
	for _, record := range tab.storage.LoadAll() {
		tab.records[record.Index] = record
	}

}

func (tab *Table[IT, DT]) Close() {
	tab.storage.Close()
}
