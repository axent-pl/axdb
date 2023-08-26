package voidstorage

import "github.com/prondos/axdb/pkg/db"

type VoidStorage[IT comparable, DT any] struct {
}

type VoidStorageMetadata struct {
}

func NewVoidStorage[IT comparable, DT any]() *VoidStorage[IT, DT] {
	return &VoidStorage[IT, DT]{}
}

func NewTable[IT comparable, DT any]() *db.Table[IT, DT] {
	storage := NewVoidStorage[IT, DT]()
	table := db.NewTable[IT, DT](storage)
	table.Open()
	return table
}
