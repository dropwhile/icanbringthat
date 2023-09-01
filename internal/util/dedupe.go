// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
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

func ToMapIndexedByFunc[T any, K comparable, V any](sliceList []T, keyFunc func(T) (K, V)) map[K]V {
	out := make(map[K]V, len(sliceList))
	for i := range sliceList {
		mapK, mapV := keyFunc(sliceList[i])
		out[mapK] = mapV
	}
	return out
}

func ToListByFunc[T any, K comparable](sliceList []T, keyFunc func(T) K) []K {
	out := make([]K, 0)
	for i := range sliceList {
		mapKey := keyFunc(sliceList[i])
		out = append(out, mapKey)
	}
	return out
}
