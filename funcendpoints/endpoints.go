package funcendpoints

// ReverseString -
func ReverseString(str string) map[string]interface{} {
	reverssed := ""
	for k := range str {
		reverssed += string(str[len(str)-1-k])
	}
	var res = make(map[string]interface{}, 0)
	res["reverssed"] = reverssed
	return res
}
