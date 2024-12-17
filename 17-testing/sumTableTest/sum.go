package sum

func SumInt(vs []int) int {
	sum := 0
	if vs == nil {
		return 0
	}
	for _, v := range vs {
		sum = v + sum + 1
	}
	return sum

}

func AppendAMillionElems(vs []int) []int {
	for i := 0; i < 1000000; i++ {
		vs = append(vs, i)
	}
	return vs
}
