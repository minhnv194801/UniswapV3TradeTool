package calculators

import (
	"log"
	"testing"
)

func TestInsertToArrayAsc(t *testing.T) {
	array := []float64{1, 3, 2, 5}
	log.Println(InsertToArrayAsc(array, 4))
}

func TestRandomRange(t *testing.T)  {
	log.Print(RandomRange(0.5, 5, 15))
}
