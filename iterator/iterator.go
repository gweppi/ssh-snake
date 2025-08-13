package iterator

type Iterator[T any] struct {
	index int
	List  []T
}

// Get next item in the list
func (i *Iterator[T]) Prev() T {
	length := len(i.List)
	// At the start of the list, so return last value in list to loop back
	if i.index == 0 {
		i.index = length - 1
		return i.List[i.index]
	}
	i.index--
	return i.List[i.index]
}

func (i *Iterator[T]) Next() T {
	length := len(i.List)
	// At the end of the list, so return first item in the list
	if i.index == length-1 {
		i.index = 0
		return i.List[i.index]
	}
	i.index++
	return i.List[i.index]
}

func (i *Iterator[T]) Current() T {
	return i.List[i.index]
}
