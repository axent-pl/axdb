package db

type Storager[IT comparable, DT any] interface {
	Store(*Record[IT, DT]) error
	Delete(Record[IT, DT]) error
	LoadAll() ([]*Record[IT, DT], error)
	Close()
}
