package util

var empty = struct{}{}

func RemoveDuplicates[T comparable](sliceList []T) []T {
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
