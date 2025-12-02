package utils

func LimitOffset[T any, S []T](s S, l, o *int) S {
	var limit int
	if l != nil {
		limit = *l
	}

	var offset int
	if o != nil {
		offset = *o
	}

	if offset > 0 {
		if offset >= len(s) {
			offset = len(s)
		}
		s = s[offset:]
	}

	if limit > 0 && len(s) > limit {
		s = s[:limit]
	}

	res := append(make(S, 0, len(s)), s...)

	return res
}
