package voidstorage

type VoidStorage[IT comparable, DT any] struct {
}

type VoidStorageMetadata struct {
}

func New[IT comparable, DT any]() *VoidStorage[IT, DT] {
	return &VoidStorage[IT, DT]{}
}
