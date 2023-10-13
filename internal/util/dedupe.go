package util

var empty = struct{}{}

func Uniq[T comparable](sliceList []T) []T {
	allKeys := make(map[T]struct{}, len(sliceList))
	list := make([]T, 0, len(sliceList))
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = empty
			list = append(list, item)
		}
	}
	return list
}

func ToSet[T comparable](sliceList []T) map[T]struct{} {
	keys := make(map[T]struct{}, len(sliceList))
	for _, item := range sliceList {
		keys[item] = empty
	}
	return keys
}

func ToSetIndexed[T comparable](sliceList []T) map[T]int {
	keys := make(map[T]int, len(sliceList))
	for i, item := range sliceList {
		keys[item] = i
	}
	return keys
}

func ToMapIndexedByFunc[T comparable, V comparable](sliceList []T, keyFunc func(T) V) map[V]T {
	out := make(map[V]T, len(sliceList))
	for i := range sliceList {
		mapKey := keyFunc(sliceList[i])
		out[mapKey] = sliceList[i]
	}
	return out
}

func ToListByFunc[T comparable, V comparable](sliceList []T, keyFunc func(T) V) []V {
	out := make([]V, 0)
	for i := range sliceList {
		mapKey := keyFunc(sliceList[i])
		out = append(out, mapKey)
	}
	return out
}
