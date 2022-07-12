package filewriter

import (
	"encoding/json"
	"testing"
)

func TestJSONWriter_WriteData(t *testing.T) {
	data := map[string]struct {
		TotalQuantity float64 `json:"total_quantity" xml:"total_quantity"`
		TotalPayable  float64 `json:"total_payable" xml:"total_payable"`
	}{}
	data["pubkey_1"] = struct {
		TotalQuantity float64 `json:"total_quantity" xml:"total_quantity"`
		TotalPayable  float64 `json:"total_payable" xml:"total_payable"`
	}{
		TotalQuantity: 0,
		TotalPayable:  0,
	}
	data["pubkey_2"] = struct {
		TotalQuantity float64 `json:"total_quantity" xml:"total_quantity"`
		TotalPayable  float64 `json:"total_payable" xml:"total_payable"`
	}{
		TotalQuantity: 0,
		TotalPayable:  0,
	}
	b, _ := json.MarshalIndent(data, "", "  ")
	err := NewJSONWriter("./data"+"/", "my_pubkey").WriteData(b)
	if err != nil {
		t.Fatal(err)
	}
}
