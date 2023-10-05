package voidstorage

import "github.com/prondos/axdb/pkg/db"

func (s *VoidStorage[IT, DT]) Store(*db.Record[IT, DT]) error {
	return nil
}

func (s *VoidStorage[IT, DT]) Delete(db.Record[IT, DT]) error {
	return nil
}

func (s *VoidStorage[IT, DT]) LoadAll() ([]*db.Record[IT, DT], error) {
	var records []*db.Record[IT, DT]
	return records, nil
}

func (s *VoidStorage[IT, DT]) Close() {

}
