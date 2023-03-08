package state

import (
	"testing"
)

func TestFormatTime(t *testing.T) {

	tests := []struct {
		tflt float32
		tstr string
	}{{269, "4:29"},
		{269.312, "4:29"},
		{3868.21, "1:4:28"},
		{0.123, "0:00"},
	}

	for _, tst := range tests {
		str := TrkTimeToString(tst.tflt)
		if str != tst.tstr {
			t.Errorf("time conversion fail: %f -> %s was %s", tst.tflt, tst.tstr, str)
		}
	}
}
