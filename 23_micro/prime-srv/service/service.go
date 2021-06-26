package service

func GetPrime(input ...int64) (ret []int64) {
	for _, v := range input {
		if isPrime(v) {
			ret = append(ret, v)
		}
	}
}

func isPrime(v int64) bool {
	var j int64 = 2
	for ; j < v/2; j++ {
		if v%j == 0 {
			return false
		}
	}
	return true
}
