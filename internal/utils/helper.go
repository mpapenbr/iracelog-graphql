package utils

func IntSliceToInt32Slice(in []int) []int32 {
	ret := make([]int32, len(in))
	for i, v := range in {
		ret[i] = int32(v)
	}
	return ret
}
