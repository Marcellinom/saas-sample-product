package common

type InfiniteScrollResult[T any] struct {
	nextCursor string
	data       []T
	total      int
}

func NewInfiniteScrollResult[T any](nextCursor string, data []T, total int) InfiniteScrollResult[T] {
	return InfiniteScrollResult[T]{nextCursor: nextCursor, data: data, total: total}
}

func (r InfiniteScrollResult[T]) NextCursor() string {
	return r.nextCursor
}

func (r InfiniteScrollResult[T]) Data() []T {
	return r.data
}

func (r InfiniteScrollResult[T]) Total() int {
	return r.total
}
