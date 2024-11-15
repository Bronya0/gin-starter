package map_util

func GetMapKeys[K comparable, V any](m map[K]V) []K {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
