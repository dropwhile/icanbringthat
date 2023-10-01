package util

var empty = struct{}{}

func RemoveDuplicates[T string | int](sliceList []T) []T {
	allKeys := make(map[T]struct{})
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = empty
			list = append(list, item)
		}
	}
	return list
}

func ToSet[T string | int](sliceList []T) map[T]struct{} {
	keys := make(map[T]struct{})
	for _, item := range sliceList {
		keys[item] = empty
	}
	return keys
}

func ToSetIndexed[T string | int](sliceList []T) map[T]int {
	keys := make(map[T]int)
	for i, item := range sliceList {
		keys[item] = i
	}
	return keys
}
