package utils

type Counter[T comparable] map[T]int

func (c Counter[T]) Add(elem T) {
	c[elem]++
}

func (c Counter[T]) Has(key T) bool {
	_, exists := c[key]
	return exists
}

func (c Counter[T]) Count(key T) int {
	value := c[key]
	return value
}

func (c Counter[T]) GT(times int) []T {
	if c == nil {
		return nil
	}

	ret := make([]T, 0, len(c))
	for val, count := range c {
		if count > times {
			ret = append(ret, val)
		}
	}

	return ret
}

func (c Counter[T]) ToSlice() []T {
	if c == nil {
		return nil
	}

	ret := make([]T, 0, len(c))
	for val := range c {
		ret = append(ret, val)
	}

	return ret
}

func ToCounter[S ~[]TArg, TKey comparable, TArg any](collection S, selector func(arg TArg) TKey) Counter[TKey] {
	result := make(Counter[TKey], len(collection))
	for _, item := range collection {
		result[selector(item)]++
	}
	return result
}
