package iterator

type Iterator[T any] struct {
	Index int
	List  []T
}

// Get next item in the list
func (i *Iterator[T]) Prev() T {
	length := len(i.List)
	// At the start of the list, so return last value in list to loop back
	if i.Index == 0 {
		i.Index = length - 1
		return i.List[i.Index]
	}
	i.Index--
	return i.List[i.Index]
}

func (i *Iterator[T]) Next() T {
	length := len(i.List)
	// At the end of the list, so return first item in the list
	if i.Index == length-1 {
		i.Index = 0
		return i.List[i.Index]
	}
	i.Index++
	return i.List[i.Index]
}

func (i *Iterator[T]) Current() T {
	return i.List[i.Index]
}
