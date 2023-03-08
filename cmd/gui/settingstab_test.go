package gui

import (
	"testing"
)

func TestIsByte(t *testing.T) {
	var valid []string = []string{"1", "0", "255", "123"}
	var invalid []string = []string{"-1", "a1", "256", "12.3"}

	for _, b := range valid {
		if isWord(b) != nil {
			t.Errorf("invalid: %s", b)
		}
	}
	for _, b := range invalid {
		if isWord(b) == nil {
			t.Errorf("invalid: %s", b)
		}
	}
}

func TestIsIPAddr(t *testing.T) {
	var valid []string = []string{"123.11.22.21", "0.0.0.0", "127.0.0.1", "192.168.0.95"}
	var invalid []string = []string{"1.2.3", "11.22.33.999", "1", "-1.23.33.44"}

	for _, b := range valid {
		if isIPaddr(b) != nil {
			t.Errorf("invalid: %s", b)
		}
	}
	for _, b := range invalid {
		if isIPaddr(b) == nil {
			t.Errorf("invalid: %s", b)
		}
	}
}
