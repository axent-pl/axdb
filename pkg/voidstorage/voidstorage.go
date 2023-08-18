package voidstorage

import "github.com/prondos/axdb/pkg/db"

type VoidStorage[IT comparable, MT any, DT any] struct {
}

type VoidStorageMetadata struct {
}

func NewVoidStorage[IT comparable, MT VoidStorageMetadata, DT any]() *VoidStorage[IT, MT, DT] {
	return &VoidStorage[IT, MT, DT]{}
}

func NewTable[IT comparable, DT any]() *db.Table[IT, VoidStorageMetadata, DT] {
	storage := NewVoidStorage[IT, VoidStorageMetadata, DT]()
	table := db.NewTable[IT, VoidStorageMetadata, DT](storage)
	table.Open()
	return table
}
