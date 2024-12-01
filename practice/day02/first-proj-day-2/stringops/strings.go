package stringops

func ReverseAndUppercase(s1, s2 string) string {
	reversed := reverseString(s1) + reverseString(s2)
	upper := toUpperCase(reversed)
	return upper
}
