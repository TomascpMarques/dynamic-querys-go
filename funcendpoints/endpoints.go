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

// ReverseStringBool -
func ReverseStringBool(bol bool, str string) map[string]interface{} {
	if !bol {
		var res = make(map[string]interface{}, 0)
		res["reverssed"] = str
		res["reversse"] = bol
		return res
	}
	reverssed := ""
	for k := range str {
		reverssed += string(str[len(str)-1-k])
	}
	var res = make(map[string]interface{}, 0)
	res["reverssed"] = reverssed
	res["reversse"] = bol
	return res
}

// TakeAInterfaceArray -
func TakeAInterfaceArray(array []interface{}) []interface{} {
	return array
}

// TakeAMap -
func TakeAMap(array map[string]interface{}) map[string]interface{} {
	return array
}
