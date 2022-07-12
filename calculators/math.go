package calculators

import (
	"math/big"
	"math/rand"
	"time"
)

func ExistedFloat64(array []float64, check float64) bool {
	for _, a := range array {
		if a > check {
			return false
		}
		if a == check {
			return true
		}
	}
	return false
}

// insert a number, randomly separate this number into a number of parts
// min percent is the smallest volume a part has
// 100%, min, part
// block = 100/(part); f = block - min;
// p = i * block
func RandomRange(num float64, part uint8, minPercent float64) []float64 {
	ret := []float64{}
	rand.Seed(time.Now().UnixNano())
	block := DivFloat64(num, float64(part))
	gapPart := SubFloat64(block, MulFloat64(DivFloat64(minPercent, 100), num))
	for i := uint8(0); i < part; i++ {
		percent := rand.Float64()
		value := MulFloat64(percent, gapPart)
		//if ExistedFloat64(ret, value) || value == 0 {
		//	i--
		//	continue
		//}

		ret = InsertToArrayAsc(ret, AddFloat64(MulFloat64(block, float64(i)), value))
	}
	bias := SubFloat64(num, ret[len(ret)-1])
	for i, v := range ret {
		ret[i] = AddFloat64(v, bias)
	}
	return ret
}

func InsertToArrayAsc(array []float64, number float64) []float64 {
	// insert
	for i, num := range array {
		if num > number {
			right := make([]float64, len(array[i:]))
			copy(right, array[i:])
			array = append(append(array[:i], number), right...)
			return array
		}
	}
	array = append(array, number)
	return array
}

func DivBy10Float64(x float64, index int64) float64 {
	for index > 0 {
		x = DivFloat64(x, 10)
		index--
	}

	return x
}

func DivFloat64(x float64, y float64) float64 {
	f, _ := new(big.Float).Quo(big.NewFloat(x), big.NewFloat(y)).Float64()
	return f
}

func MulFloat64(x float64, y float64) float64 {
	f, _ := new(big.Float).Mul(big.NewFloat(x), big.NewFloat(y)).Float64()
	return f
}

func AddFloat64(minus float64, subtrahend float64) float64 {
	f, _ := new(big.Float).Add(big.NewFloat(minus), big.NewFloat(subtrahend)).Float64()
	return f
}

func SubFloat64(minus float64, subtrahend float64) float64 {
	f, _ := new(big.Float).Sub(big.NewFloat(minus), big.NewFloat(subtrahend)).Float64()
	return f
}
