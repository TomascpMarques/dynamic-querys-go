package functions

// Test -
func Test(num1, num2, num3 int) int {
	//fmt.Println(num1 + num2 + num3)
	return num1 + num2 + num3
}

// Test2 -
func Test2(num1, num2, num3 int) int {
	//fmt.Println(num1 * num2 * num3)
	return num1 * num2 * num3
}

// Test3 -
func Test3(num1, num2, num3 float64) (float64, float64) {
	//fmt.Println(num1 * num2 * num3)
	return num1 * num2 * num3, num1 * num2
}

// Test4 -
func Test4(fields []string, id string, i float64) ([]string, string, float64) {
	return fields, id, i
}
