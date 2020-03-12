package convert

import (
	"strconv"
)


func StringToBytes(str string) []byte {
	return []byte(str)
}

func BytesToString(byteData []byte) string {
	return string(byteData)
}

//字符串转in 适用于10进制
func StringToInt(str string) int {
	intResult, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		intResult = 0
	}

	return int(intResult)
}

func StringToInt32(str string) int32 {
	intResult, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		intResult = 0
	}

	return int32(intResult)
}

//字符串转in64 适用于10进制
func StringToInt64(str string) int64 {
	intResult, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		intResult = 0
	}

	return intResult
}

//字符串转float64
func StringToFloat64(str string) float64 {
	floatResult, err := strconv.ParseFloat(str, 64)
	if err != nil {
		floatResult = 0
	}
	return floatResult
}

//整数转字符串 支持int int8 int16 int32 int64 uint uint8 uint16 uint32 uint64
func IntToString(intData int) string {
	str := strconv.FormatInt(int64(intData), 10)
	return str
}

func Int8ToString(intData int8) string {
	str := strconv.FormatInt(int64(intData), 10)
	return str
}

func Int16ToString(intData int16) string {
	str := strconv.FormatInt(int64(intData), 10)
	return str
}

func Int32ToString(intData int32) string {
	str := strconv.FormatInt(int64(intData), 10)
	return str
}

func Int64ToString(intData int64) string {
	str := strconv.FormatInt(int64(intData), 10)
	return str
}

func UintToString(intData uint) string {
	str := strconv.FormatUint(uint64(intData), 10)
	return str
}

func Uint8ToString(intData uint8) string {
	str := strconv.FormatUint(uint64(intData), 10)
	return str
}

func Uint16ToString(intData uint16) string {
	str := strconv.FormatUint(uint64(intData), 10)
	return str
}

func Uint32ToString(intData uint32) string {
	str := strconv.FormatUint(uint64(intData), 10)
	return str
}

func Uint64ToString(intData uint64) string {
	str := strconv.FormatUint(uint64(intData), 10)
	return str
}

//注意精度损失
func Float64ToString(floatData float64) string {
	return strconv.FormatFloat(floatData, 'f', -1, 64)
}

//注意精度损失
func Float32ToString(floatData float32) string {
	return strconv.FormatFloat(float64(floatData), 'f', -1, 64)
}
