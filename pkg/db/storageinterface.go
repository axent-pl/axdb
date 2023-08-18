package db

type StorageInterface[IT comparable, MT any, DT any] interface {
	Store(*Record[IT, MT, DT]) error
	Delete(Record[IT, MT, DT]) error
	LoadAll() []*Record[IT, MT, DT]
	Close()
}
