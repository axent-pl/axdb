package db

type Record[IT comparable, DT any] struct {
	Index IT
	Data  DT
}

func NewRecord[IT comparable, DT any](index IT, data DT) *Record[IT, DT] {
	rec := &Record[IT, DT]{
		Index: index,
		Data:  data,
	}
	return rec
}
