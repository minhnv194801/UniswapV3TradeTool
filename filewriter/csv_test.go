package filewriter

import (
	"testing"
)

func TestCSVWriter_WriteData(t *testing.T) {
	writer := NewCSVWriter("./data", "my_pubkey", []string{"a", "b", "c", "d", "e"})
	err := writer.WriteData([]string{"abc", "def", "ghi", "jkl", "mno"})
	if err != nil {
		t.Fatal(err)
	}
}
