package voidstorage

import "github.com/prondos/axdb/pkg/db"

func (s *VoidStorage[IT, MT, DT]) Store(*db.Record[IT, MT, DT]) error {
	return nil
}

func (s *VoidStorage[IT, MT, DT]) Delete(db.Record[IT, MT, DT]) error {
	return nil
}

func (s *VoidStorage[IT, MT, DT]) LoadAll() []*db.Record[IT, MT, DT] {
	var records []*db.Record[IT, MT, DT]
	return records
}

func (s *VoidStorage[IT, MT, DT]) Close() {

}
