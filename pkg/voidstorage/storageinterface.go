package voidstorage

import "github.com/prondos/axdb/pkg/db"

func (s *VoidStorage[IT, DT]) Store(*db.Record[IT, DT]) error {
	return nil
}

func (s *VoidStorage[IT, DT]) Delete(db.Record[IT, DT]) error {
	return nil
}

func (s *VoidStorage[IT, DT]) LoadAll() []*db.Record[IT, DT] {
	var records []*db.Record[IT, DT]
	return records
}

func (s *VoidStorage[IT, DT]) Close() {

}
