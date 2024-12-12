package cachestorage

type CacheStorageFuncHandler[T any] interface {
	SetFunc(func(int) bool)
	GetFunc() func(int) bool
	Comparison(CacheStorageFuncHandler[T]) bool
}
