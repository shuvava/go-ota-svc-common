package data_test

import (
	"fmt"
	"testing"

	"github.com/shuvava/go-ota-svc-common/data"
)

func TestValidHex(t *testing.T) {
	cases := []struct {
		Hex      string
		Len      int
		Expected bool
	}{
		{Hex: "ab12", Len: 4, Expected: true},
		{Hex: "abq12", Len: 5, Expected: false},
		{Hex: "/ab12", Len: 5, Expected: false},
		{Hex: "ab:12", Len: 5, Expected: false},
		{Hex: "ab1?2", Len: 5, Expected: false},
		{Hex: "ab12[", Len: 5, Expected: false},
		{Hex: "Ab12", Len: 4, Expected: false},
		{Hex: "Ab12", Len: 2, Expected: false},
		{Hex: "0123456789abcdef", Len: 16, Expected: true},
	}
	for _, test := range cases {
		exStr := "valid"
		if test.Expected != true {
			exStr = "invalid"
		}
		name := fmt.Sprintf("string '%s' with len %d is %s", test.Hex, test.Len, exStr)
		t.Run(name, func(t *testing.T) {
			got := data.ValidHex(test.Len, test.Hex)
			if got != test.Expected {
				t.Errorf("got %v, want %v", got, test.Expected)
			}
		})
	}
}

func TestValidBase64(t *testing.T) {
	cases := []struct {
		Base64   string
		Expected bool
	}{
		{Base64: "WvLTlMrX9NpYDQlEIFlnDA==", Expected: true},
		{Base64: "Invalid string", Expected: false},
	}
	for _, test := range cases {
		exStr := "valid"
		if test.Expected != true {
			exStr = "invalid"
		}
		name := fmt.Sprintf("string '%s' is %s", test.Base64, exStr)
		t.Run(name, func(t *testing.T) {
			got := data.ValidBase64(test.Base64)
			if got != test.Expected {
				t.Errorf("got %v, want %v", got, test.Expected)
			}
		})
	}
}
