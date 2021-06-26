package service

func GetSum(inputs ...int64) (res int64) {
	for _, v := range inputs {
		res += v
	}
	return res
}
