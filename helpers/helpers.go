package helpers

import (
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

func Absolute(num int) int {
	return int(math.Abs(float64(num)))
}

func GetSign(v int) int {
	if v == 0 {
		return 0
	}
	return v / Absolute(v)
}

func StringToInt(s string) int {
	num, _ := strconv.Atoi(s)
	return num
}

func StringToInt64(s string) uint64 {
	num, _ := strconv.Atoi(s)
	return uint64(num)
}

func ReadFile() []string {
	file, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	return strings.Split(string(file), "\n")
}
