package utils

type GroupBy[TKey comparable, TArg any] map[TKey][]TArg

func (m GroupBy[TKey, TArg]) Add(key TKey, value TArg) {
	m[key] = append(m[key], value)
}

func (m GroupBy[TKey, TArg]) Delete(key TKey) {
	delete(m, key)
}

func (m GroupBy[TKey, TArg]) Keys() []TKey {
	set := Set[TKey]{}
	for k := range m {
		set.Add(k)
	}
	return set.ToSlice()
}

func (m GroupBy[TKey, TArg]) Has(key TKey) bool {
	_, exists := m[key]
	return exists
}

func ToGroupBy[S ~[]TArg, TKey comparable, TValue, TArg any](collection S, selector func(arg TArg) (TKey, TValue)) GroupBy[TKey, TValue] {
	result := make(GroupBy[TKey, TValue], len(collection))
	for _, element := range collection {
		k, v := selector(element)
		result[k] = append(result[k], v)
	}
	return result
}
