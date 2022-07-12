package utils

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func Parse18Decimal(s string) string {
	remainString := s
	if len(remainString) < 18 {
		remainString = fmt.Sprintf("%018d", 0) + remainString
	}
	remainString = remainString[:len(remainString)-18] + "." + remainString[len(remainString)-18:]
	return remainString
}

func ParseStringTo18Decimal(s string) *big.Int {
	splits := strings.Split(s, ".")
	if len(splits) == 1 {
		s += fmt.Sprintf("%018d", 0)
	} else {
		left := splits[0]
		right := splits[1]
		if len(right) < 18 {
			right += fmt.Sprintf("%018d", 0)
		}
		s = left + right[:18]
	}
	bigInt, _ := new(big.Int).SetString(s, 10)
	return bigInt
}

func Parse18DecimalToFloat(s string) (float64, error) {
	remainString := s
	if len(remainString) < 18 {
		remainString = fmt.Sprintf("%018d", 0) + remainString
	}
	remainString = remainString[:len(remainString)-18] + "." + remainString[len(remainString)-18:]
	return strconv.ParseFloat(remainString, 64)
}
