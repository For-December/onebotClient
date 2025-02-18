package onebotClient

import "sort"

// IsTargetInArray 二分查找
func IsTargetInArray[T int | float64 | string | int64](target T, array []T) bool {
	// 切片必须升序
	sort.Slice(array, func(i, j int) bool {
		return array[i] < array[j]
	})
	//index := sort.SearchStrings(array, target)
	index := sort.Search(len(array), func(i int) bool { return array[i] >= target })
	//index的取值：0 ~ (len(str_array)-1)
	return index < len(array) && array[index] == target
}
