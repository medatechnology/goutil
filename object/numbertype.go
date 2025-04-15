package object

// Define helper for number declaration/interface
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
type Integer interface {
	Signed | Unsigned
}
type Float interface {
	~float32 | ~float64
}
type Number interface {
	Integer | Float
}

// Use var so later we can change it, it is crude for now
// but can do: object.THOUSAND_SEPARATOR='.' before calling
// any of the number methods
var (
	THOUSAND_SEPARATOR byte = ','
	DECIMAL_SEPARATOR  byte = '.'
)

// Return the maximum, Minimum values and index in the array of integer
func MaxMinIntArray(arr []int) (max, min, maxindex, minindex int) {
	max = -1
	min = 2147483647
	maxindex = -1
	minindex = -1

	for i, v := range arr {
		if max < v {
			max = v
			maxindex = i
		}
		if min > v {
			min = v
			minindex = i
		}
	}
	return
}

// Return true if num exists in arr
func ExistInArrayInt(num int, arr []int) bool {
	for _, r := range arr {
		if r == num {
			return true
		}
	}
	return false
}

// Simpler generic ABS func, no need to import math. Need Number Type
func Abs[T Number](num T) T {
	if num < 0 {
		return -num
	}
	return num
}

// NOTE: not needed to check the type/kind for unsigned, just negate
// func Abs[T Number](number T) T {
// 	val := reflect.ValueOf(number)
// 	switch val.Kind() {
// 	// case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
// 	// case reflect.Float32, reflect.Float64:
// 	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr:
// 					return number
// 	}
// 	if number < 0 {
// 					return -number
// 	}
// 	return T(number)
// }

// PLEASE NOTE TO USE THIS ONLY IF YOU ARE SURE THAT THE STRING IS INTEGER
// Mostly only to convert simple string index back to positive integer
// NOTE: only works on decimals and POSITIVE value, so string "-123" will not work
func IntPlus(str string, zeroValue int) int {
	if len(str) < 1 {
		return zeroValue
	}
	var value int
	for _, char := range str {
		value = value*10 + int(char-'0')
	}
	return value
}

// Faster Atoi, only for string number that is less than 10 digits.
// No error checking, if there is error, then return 0.
// Skipping the decimals (anything after . like: 123.123 ==> 123)
// Ignoring the thousand separator
// abs = true, return only absolute value.
// NOTE: only works on round numbers (integer) cannot be decimals
func Int(smallNumber string, abs bool) int {
	s0 := smallNumber
	if smallNumber[0] == '-' || smallNumber[0] == '+' {
		if len(smallNumber) < 2 {
			return 0
		}
		smallNumber = smallNumber[1:]
	}
	n := 0
	for _, ch := range []byte(smallNumber) {
		// Skip thousand separator
		if ch == THOUSAND_SEPARATOR {
			continue
		}
		// Skip the decimals
		if ch == DECIMAL_SEPARATOR {
			break
		}
		ch -= '0'
		if ch > 9 {
			return 0
		}
		n = n*10 + int(ch)
	}
	if !abs {
		if s0[0] == '-' {
			n = -n
		}
	}
	return n
}
