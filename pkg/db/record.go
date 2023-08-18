package db

type Record[IT comparable, MT any, DT any] struct {
	Index    IT
	Data     DT
	Metadata *MT
}

func NewRecord[IT comparable, MT any, DT any](index IT, data DT) *Record[IT, MT, DT] {
	rec := &Record[IT, MT, DT]{
		Index: index,
		Data:  data,
	}
	return rec
}
