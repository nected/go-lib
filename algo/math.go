package algo

func MAX(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MIN(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func In[Key comparable](key Key, list []Key) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

func Intersect[Key string | int | bool | uint | float32](list1 []Key, list2 []Key) (finalList []Key) {
	for _, v := range list1 {
		if In(v, list2) {
			finalList = append(finalList, v)
		}
	}
	return finalList
}

func Union[Key string | int | bool | uint | float32](list1 []Key, list2 []Key) (finalList []Key) {
	m := map[Key]bool{}
	for _, v := range list1 {
		m[v] = true
	}

	for _, v := range list2 {
		m[v] = true
	}

	for k := range m {
		finalList = append(finalList, k)
	}
	return finalList
}

func AminusB[Key string | int | bool | uint | float32](list1 []Key, list2 []Key) (finalList []Key) {
	for _, k := range list1 {
		if !In(k, list2) {
			finalList = append(finalList, k)
		}
	}
	return finalList
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
