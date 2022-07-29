package utils


//GetIndex 获取 ele 在 slice 中的索引
func GetIndex(sliceArr []string, ele string) int {

	for index, value := range sliceArr {
		if value == ele {
			return index
		}
	}

	return -1
}

//Difference ----求 slice1-(slice1,slice2的交集)后的集合
func Difference(slice1, slice2 []string) []string {

	res := make([]string, 0)
	tmpMap := make(map[string]int)

	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		tmpMap[v]++
	}

	for _, value := range slice1 {
		num, _ := tmpMap[value]
		if num == 0 {
			res = append(res, value)
		}
	}
	return res
}

//Intersect ----求 slice1, slice2的交集
func Intersect(slice1, slice2 []string) []string {

	res := make([]string, 0)
	tmpMap := make(map[string]int)
	for _, v := range slice1 {
		tmpMap[v]++
	}

	for _, v := range slice2 {
		if _, ok := tmpMap[v]; ok {
			res = append(res, v)
		}
	}

	return res
}

//数组去重
func RemoveDuplicationByStringSlice(arr []string) []string {

	tmpMap := make(map[string]int)
	i := 0
	for _, ele := range arr {
		if _, ok := tmpMap[ele]; ok {
			continue
		}

		tmpMap[ele] = 1
		arr[i] = ele
		i++
	}

	return arr[:i]
}

//AddElementAfterIndex ----add element after index=i
func AddElementAfterIndex(i int, slice, newSlice []string) []string {

	switch i {
	case 0:
		slice = append(slice[:1], append(newSlice, slice[1:]...)...)
	case len(slice) - 1:
		slice = append(slice, newSlice...)
	default:
		slice = append(slice[:i+1], append(newSlice, slice[i+1:]...)...)
	}

	return slice
}

//AddElementBeforeIndex ----add element before index=i
func AddElementBeforeIndex(i int, slice, newSlice []string) []string {

	switch i {
	case 0:
		slice = append(newSlice, slice...)
	default:
		slice = append(slice[:i], append(newSlice, slice[i:]...)...)
	}

	return slice
}

//ReplaceElement ----use newSlice replace element at index=i
func ReplaceElement(i int, slice, newSlice []string) []string {

	switch i {
	case 0:
		slice = append(newSlice, slice[1:]...)
	case len(slice) - 1:
		slice = append(slice[:len(slice)-1], newSlice...)
	default:
		slice = append(slice[:i], append(newSlice, slice[i+1:]...)...)
	}

	return slice
}

func DeleteElement(i int, slice []string) []string {

	switch i {
	case 0:
		slice = slice[1:]
	case len(slice) - 1:
		slice = slice[:len(slice)-1]
	default:
		slice = append(slice[:i], slice[i+1:]...)
	}

	return slice
}

